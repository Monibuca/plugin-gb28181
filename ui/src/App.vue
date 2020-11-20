<template>
    <div>
        <div class="tabpanel" v-if="$parent.titleTabActive === 0">
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
                                <mu-button flat v-if="item.Connected" @click="ptz(prop.row.ID, $index,item)">云台
                                </mu-button>
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
        </div>
        <div class="tabpanel" v-if="$parent.titleTabActive === 1">
            <div class="search">
            <i-select v-model="channelSelectedList" multiple>
              <Option v-for="(channel,index) in channelList"
                      :value="channel"
                      :label="channel.DeviceID + '(' + channel.Name + ')'"
                      :key="index">
              </Option>
            </i-select>
            </div>
            <div class="flex-box">
                <template v-for="(channel,index) in channelSelectedList">
                    <div class="flex-item" :key="index" v-if="channel.ID">
                        <webrtc-player2 :stream-path="'gb28181/'+channel.ID"></webrtc-player2>
                    </div>
                </template>
            </div>
        </div>
        <webrtc-player ref="player" @ptz="sendPtz" v-model="previewStreamPath" :PublicIP="PublicIP"></webrtc-player>
    </div>
</template>
<script>
    import WebrtcPlayer from "./components/Player"
    import WebrtcPlayer2 from "./components/Player2";
    import {getPTZCmd, PTZ_TYPE} from "./utils/ptz-cmd";

    export default {
        components: {
            WebrtcPlayer,
            WebrtcPlayer2
        },
        props: {
            ListenAddr: String
        },
        computed: {
            PublicIP() {
                return this.ListenAddr.split(":")[0]
            }
        },
        data() {
            return {
                Devices: [], previewStreamPath: false,
                channelList: [],
                channelSelectedList:[],
                context: {
                    id: null,
                    channel: 0,
                    item: null
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
            };
        },
        created() {
            this.fetchlist();
        },
        mounted() {
            this.$parent.titleTabs = ["列表", "N路播放"];
        },
        methods: {
            fetchlist() {
                const listES = new EventSource(this.apiHost + "/gb28181/list");
                listES.onmessage = (evt) => {
                    if (!evt.data) return;
                    this.Devices = JSON.parse(evt.data) || [];
                    this.Devices.sort((a, b) => (a.ID > b.ID ? 1 : -1));
                    let channelList =[]
                    this.Devices.forEach((device)=>{
                      const channels = device.Channels || [];
                      if(channels.length > 0){
                          channelList = channelList.concat(channels);
                      }
                    });
                    if(channelList.length > 0){
                        this.channelList = channelList;
                    }
                };
                this.$once("hook:destroyed", () => listES.close());
            },
            ptz(id, channel, item) {
                this.context = {
                    id, channel, item
                };
                this.previewStreamPath = true
                this.$nextTick(() => this.$refs.player.play("gb28181/" + item.DeviceID));
            },
            sendPtz(options) {
                const ptzCmd = getPTZCmd(options);
                const ptzCmdStop = getPTZCmd({type: PTZ_TYPE.stop});
                this.ajax.get("/gb28181/control", {
                    id: this.context.id,
                    channel: this.context.channel,
                    ptzcmd: ptzCmd,
                }).then(x => {
                    if (options.type === PTZ_TYPE.stop || options.cycle === true) {
                        return;
                    }
                    setTimeout(() => {
                        this.ajax.get("/gb28181/control", {
                            id: this.context.id,
                            channel: this.context.channel,
                            ptzcmd: ptzCmdStop,
                        });
                    }, 500)
                });
            },
            invite(id, channel, item) {
                this.ajax.get("/gb28181/invite", {id, channel}).then(x => {
                    item.Connected = true
                });
            },
            bye(id, channel, item) {
                this.ajax.get("/gb28181/bye", {id, channel}).then(x => {
                    item.Connected = false
                });
            }
        },
    };
</script>
<style scoped>
    .flex-box {
        display: flex;
        flex-flow: row wrap;
        align-content: flex-start;
    }

    .flex-item {
        flex: 0 0 33.3333%;
        height: 275px;
        padding: 10px;
    }
</style>
