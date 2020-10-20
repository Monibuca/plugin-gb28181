<template>
  <div>
    <mu-data-table :data="Devices" :columns="columns">
      <template #expand="prop">
        <mu-data-table :data="prop.row.Channels" :columns="columns2">
          <template #default="{ row: item, $index }">
            <td>{{ item.DeviceID }}</td>
            <td>{{ item.Name }}</td>
            <td>{{ item.Manufacturer }}</td>
            <td>{{ item.Address }}</td>
            <td>{{ item.Status }}</td>
            <td>
              <mu-button flat v-if="item.Connected" @click="ptz(prop.row.ID, $index,item)">云台</mu-button>
              <mu-button flat v-if="item.Connected" @click="bye(prop.row.ID, $index)">断开</mu-button>
              <mu-button v-else flat @click="invite(prop.row.ID, $index,item)"
              >连接
              </mu-button
              >
            </td>
          </template>
        </mu-data-table>
      </template>
      <template #default="{ row: item }">
        <td>{{ item.ID }}</td>
        <td>{{ item.Channels ? item.Channels.length : 0 }}</td>
        <td>
          <StartTime :value="item.RegisterTime"></StartTime>
        </td>
        <td>
          <StartTime :value="item.UpdateTime"></StartTime>
        </td>
        <td>{{ item.Status }}</td>
      </template>
    </mu-data-table>
    <webrtc-player ref="player" @ptz="sendPtz" v-model="previewStreamPath" :PublicIP="PublicIP"></webrtc-player>
  </div>
</template>
<script>
import WebrtcPlayer from "./components/Player"
export default {
  components:{
    WebrtcPlayer
  },
  props:{
    ListenAddr:String
  },
  computed:{
    PublicIP(){
      return this.ListenAddr.split(":")[0]
    }
  },
  data() {
    return {
      Devices: [], previewStreamPath:false,
      context:{
        id:null,
        channel:0,
        item:null
      },
      columns: Object.freeze(
          ["设备号", "通道数", "注册时间", "更新时间", "状态"].map(
              (title) => ({
                title,
              })
          )
      ),
      columns2: Object.freeze([
        "通道编号",
        "名称",
        "厂商",
        "地址",
        "状态",
        "操作",
      ]).map((title) => ({title})),
      ptzCmds:["A50F010800880045","A50F01018800003E", "A50F010400880041","A50F01028800003F"]
    };
  },
  created() {
    this.fetchlist();
  },
  methods: {
    fetchlist() {
      const listES = new EventSource(this.apiHost + "/gb28181/list");
      listES.onmessage = (evt) => {
        if (!evt.data) return;
        this.Devices = JSON.parse(evt.data) || [];
        this.Devices.sort((a, b) => (a.ID > b.ID ? 1 : -1));
      };
      this.$once("hook:destroyed", () => listES.close());
    },
    ptz(id, channel,item) {
      this.context = {
        id,channel,item
      }
      this.previewStreamPath = true
      this.$nextTick(() =>this.$refs.player.play("gb28181/"+item.DeviceID));
    },
    sendPtz(n){
      this.ajax.get("/gb28181/control", {
        id:this.context.id,
        channel:this.context.channel,
        ptzcmd: this.ptzCmds[n-1],
      }).then(x=>{
        setTimeout(()=>{
          this.ajax.get("/gb28181/control", {
            id:this.context.id,
            channel:this.context.channel,
            ptzcmd: "A50F0100000000B5",
          });
        },1000)
      });
    },
    invite(id, channel,item) {
      this.ajax.get("/gb28181/invite", {id, channel}).then(x=>{
        item.Connected = true
      });
    },
    bye(id, channel,item) {
      this.ajax.get("/gb28181/bye", {id, channel}).then(x=>{
        item.Connected = false
      });;
    }
  },
};
</script>