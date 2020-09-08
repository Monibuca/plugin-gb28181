<template>
  <div>
    <mu-data-table :data="Devices" :columns="columns">
      <template #default="{row:item}">
        <td>{{item.ID}}</td>
        <td>
          <StartTime :value="item.RegisterTime"></StartTime>
        </td>
        <td>
          <StartTime :value="item.UpdateTime"></StartTime>
        </td>
        <td>{{item.IP}}</td>
        <td>{{item.Status}}</td>
        <td>
          <mu-button flat @click="preview(item)">预览</mu-button>
        </td>
      </template>
    </mu-data-table>
  </div>
</template>
<script>
export default {
  data(){
    return {
      Devices:[],
      columns:Object.freeze([
          "设备号","注册时间","更新时间","IP","状态"
      ].map(title=>({title})))
    }
  },
  created(){
    this.fetchlist()
  },
  methods:{
    fetchlist() {
      const listES = new EventSource(this.apiHost + "/gb28181/list");
      listES.onmessage = evt => {
        if (!evt.data) return;
        this.Devices = JSON.parse(evt.data) || [];
      };
      this.$once("hook:destroyed",()=>listES.close())
    },
    preview(item){
      alert("尚未实现")
    }
  }
}
</script>