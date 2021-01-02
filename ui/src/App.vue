<template>
    <div>
        <div class="tabpanel" v-if="$parent.titleTabActive === 0">
            <mu-data-table :data="Devices" :columns="columns">
                <template #expand="prop">
                    <mu-data-table
                        :data="prop.row.Channels"
                        :columns="columns2"
                    >
                        <template #default="{ row: item, $index }">
                            <td>{{ item.DeviceID }}</td>
                            <td>{{ item.Name }}</td>
                            <td>{{ item.Manufacturer }}</td>
                            <td>{{ item.Address }}</td>
                            <td>{{ item.Status }}</td>
                            <td>
                                <mu-button
                                    flat
                                    v-if="item.Connected"
                                    @click="ptz(prop.row.ID, $index, item)"
                                    >云台
                                </mu-button>
                                <mu-button
                                    flat
                                    v-if="item.Connected"
                                    @click="bye(prop.row.ID, $index, item)"
                                    >断开
                                </mu-button>
                                <mu-button
                                    v-else
                                    flat
                                    @click="invite(prop.row.ID, $index, item)"
                                    >连接
                                </mu-button>
                                <mu-button
                                    flat
                                    @click="
                                        getRecords(prop.row.ID, $index, item)
                                    "
                                    >录像</mu-button
                                >
                            </td>
                        </template>
                    </mu-data-table>
                </template>
                <template #default="{ row: item }">
                    <td>{{ item.ID }}</td>
                    <td>{{ item.Addr }}</td>
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
            <div class="flex-box">
                <template v-for="(channel, index) in channelShowList">
                    <div class="flex-item" :key="index" v-if="channel.DeviceID">
                        <webrtc-player2
                            :stream-path="'gb28181/' + channel.DeviceID"
                        ></webrtc-player2>
                    </div>
                </template>
            </div>
            <template v-if="channelList.length > 0">
                <Page
                    :total="channelList.length"
                    :page-size="pageInfo.onePageSize"
                    @on-change="handlePageChange"
                ></Page>
            </template>
        </div>
        <webrtc-player
            ref="player"
            @ptz="sendPtz"
            v-model="previewStreamPath"
        ></webrtc-player>
        <records
            ref="records"
            v-model="recordModal"
            :search="recordSearch"
            :channel="currentChannel"
            @close="initRecordSearch"
        ></records>
    </div>
</template>
<script>
import WebrtcPlayer from "./components/Player";
import WebrtcPlayer2 from "./components/Player2";
import Records from "./components/Records";
import { getPTZCmd, PTZ_TYPE } from "./utils/ptz-cmd";
import { getOneTimeRange } from "./utils";

export default {
    components: {
        WebrtcPlayer,
        WebrtcPlayer2,
        Records,
    },
    props: {
        ListenAddr: String,
    },
    computed: {
        PublicIP() {
            return this.ListenAddr.split(":")[0];
        },
    },
    data() {
        return {
            Devices: [],
            previewStreamPath: false,
            channelList: [],
            channelShowList: [],
            pageInfo: {
                onePageSize: 9,
                totalPage: 0,
                currentPage: 0,
            },
            currentChannel: null,
            recordModal: false,
            recordSearch: {
                id: null,
                channel: null,
                deviceId: null,
            },
            context: {
                id: null,
                channel: 0,
                item: null,
            },
            columns: Object.freeze(
                [
                    "设备号",
                    "地址",
                    "通道数",
                    "注册时间",
                    "更新时间",
                    "状态",
                ].map((title) => ({
                    title,
                }))
            ),
            columns2: Object.freeze([
                "通道编号",
                "名称",
                "厂商",
                "地址",
                "状态",
                "操作",
            ]).map((title) => ({ title })),
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
            const listES = new EventSource("/gb28181/list");
            listES.onmessage = (evt) => {
                if (!evt.data) return;
                this.Devices = JSON.parse(evt.data) || [];
                this.Devices.sort((a, b) => (a.ID > b.ID ? 1 : -1));
                let channelList = [];
                this.Devices.forEach((device) => {
                    const channels = device.Channels || [];
                    if (channels.length > 0) {
                        channelList = channelList.concat(channels);
                    }
                    if (this.recordSearch.id && this.recordSearch.deviceId) {
                        const channel = channels.find((i) => {
                            return (
                                i.DeviceID === this.recordSearch.deviceId &&
                                this.recordSearch.id === device.ID
                            );
                        });
                        if (channel) this.currentChannel = channel;
                    }
                });
                if (channelList.length > 0) {
                    this.channelList = channelList;
                    this.updatePageInfo(channelList.length);
                }
            };
            this.$once("hook:destroyed", () => listES.close());
        },
        ptz(id, channel, item) {
            this.context = {
                id,
                channel,
                item,
            };
            this.previewStreamPath = true;
            this.$nextTick(() => this.$refs.player.play(item.LiveSP));
        },
        sendPtz(options) {
            const ptzCmd = getPTZCmd(options);
            const ptzCmdStop = getPTZCmd({ type: PTZ_TYPE.stop });
            this.ajax
                .get("/gb28181/control", {
                    id: this.context.id,
                    channel: this.context.channel,
                    ptzcmd: ptzCmd,
                })
                .then((x) => {
                    if (
                        options.type === PTZ_TYPE.stop ||
                        options.cycle === true
                    ) {
                        return;
                    }
                    setTimeout(() => {
                        this.ajax.get("/gb28181/control", {
                            id: this.context.id,
                            channel: this.context.channel,
                            ptzcmd: ptzCmdStop,
                        });
                    }, 500);
                });
        },

        sendQueryRecords(options) {},

        handlePageChange(page) {
            let showList = [];
            const onePageSize = this.pageInfo.onePageSize;
            const firstIndex = page * onePageSize - onePageSize;
            const lastIndex = page * onePageSize - 1;
            showList = this.channelList.filter((item, index) => {
                return index >= firstIndex && index <= lastIndex;
            });

            this.channelShowList = showList;
            if (showList.length > 0) {
                this.pageInfo.currentPage = page;
            }
        },
        updatePageInfo(totalSize) {
            const onePageSize = this.pageInfo.onePageSize;
            let totalPage = totalSize / onePageSize;

            if (totalSize % onePageSize > 0) {
                totalPage = totalPage + 1;
            }
            this.pageInfo.totalPage = totalPage;
            if (this.pageInfo.currentPage === 0) {
                this.handlePageChange(1);
            }
        },
        invite(id, channel, item) {
            this.ajax.get("/gb28181/invite", { id, channel }).then((x) => {
                item.Connected = true;
            });
        },
        bye(id, channel, item) {
            this.ajax.get("/gb28181/bye", { id, channel }).then((x) => {
                item.Connected = false;
            });
        },

        getRecords(id, channel, item) {
            this.recordSearch.id = id;
            this.recordSearch.channel = channel;
            this.recordSearch.deviceId = item.DeviceID;
            this.recordModal = true;
        },

        initRecordSearch() {
            this.recordModal = false;
            this.recordSearch.id = null;
            this.recordSearch.channel = null;
            this.recordSearch.deviceId = null;
        },
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
    box-sizing: border-box;
    padding: 10px;
}
</style>
