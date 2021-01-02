<template>
    <Modal
        v-bind="$attrs"
        draggable
        width="900"
        v-on="$listeners"
        title="录像列表"
        @on-ok="$emit('close')"
    >
        <webrtc-player2
            v-if="channel && channel.RecordSP && player"
            :streamPath="channel.RecordSP"
        ></webrtc-player2>
        <div class="container" v-else-if="!player">
            <div class="search">
                <DatePicker
                    type="date"
                    :options="timeOptions"
                    :value="search.time"
                    placeholder="请选择时间"
                    style="width: 200px"
                    :clearable="false"
                    @on-change="handleTimeChange"
                ></DatePicker>
            </div>
            <div>
                <mu-data-table :columns="columns" :data="recordList">
                    <template #default="scope">
                        <td>{{ scope.row.DeviceID }}</td>
                        <td>{{ scope.row.Name }}</td>
                        <td>{{ scope.row.startTime }}</td>
                        <td>{{ scope.row.endTime }}</td>
                        <td>{{ scope.row.length }}</td>
                        <td>
                            <m-button @click="play(scope.row)">播放</m-button>
                        </td>
                    </template>
                </mu-data-table>
            </div>
        </div>
        <div v-else>正在连接，请稍后</div>
        <div slot="footer" v-if="player">
            <mu-button @click="back">返回</mu-button>
        </div>
    </Modal>
</template>

<script>
import { getOneTimeRange, formatTimeTips, parseTime, isDef } from "../utils";
import WebrtcPlayer2 from "./Player2";
const _now = new Date();

export default {
    name: "Records",
    components: {
        WebrtcPlayer2,
    },
    props: ["search", "channel"],
    data() {
        return {
            player: false,
            timeOptions: {
                disabledDate(date) {
                    return date && date.valueOf() > Date.now();
                },
            },
            columns: Object.freeze(
                ["设备ID", "名称", "开始时间", "结束时间", "时长", "操作"].map(
                    (title) => ({
                        title,
                    })
                )
            ),
        };
    },
    computed: {
        records() {
            return (this.channel && this.channel.Records) || [];
        },
        startTime() {
            if (!this.search.time) {
                return "";
            }
            const start = getOneTimeRange(this.search.time).start;
            const isoString = new Date(start).toISOString();
            return isoString.replace(".000Z", "");
        },
        endTime() {
            if (!this.search.time) {
                return "";
            }
            const end = getOneTimeRange(this.search.time).end;
            const isoString = new Date(end).toISOString();
            return isoString.replace(".000Z", "");
        },
        recordList() {
            const list = this.records.map((record) => {
                const startTime = new Date(record.StartTime).getTime();
                const endTime = new Date(record.EndTime).getTime();
                const timestamp = endTime - startTime;
                const timeLength = formatTimeTips(timestamp / 1000);

                const _startTime = parseTime(startTime);
                const _endTime = parseTime(endTime);
                record._startTime = (startTime / 1000) >> 0;
                record._endTime = (endTime / 1000) >> 0;
                record.length = timeLength;
                record.startTime = _startTime;
                record.endTime = _endTime;
                return record;
            });
            return list;
        },
    },
    mounted() {
        this._fetchList();
    },
    methods: {
        _fetchList() {
            if (
                isDef(this.search.id) &&
                isDef(this.search.channel) &&
                this.startTime &&
                this.endTime
            ) {
                const query = {
                    id: this.search.id,
                    channel: this.search.channel,
                    startTime: this.startTime,
                    endTime: this.endTime,
                };

                this.ajax.get("/gb28181/query/records", query).then((x) => {});
            }
        },
        handleTimeChange(date) {
            this.search.time = new Date(date);
            this._fetchList();
        },
        play(record) {
            const query = {
                id: this.search.id,
                channel: this.search.channel,
                startTime: record._startTime,
                endTime: record._endTime,
            };

            this.ajax.get("/gb28181/invite", query).then((x) => {});
            this.player = true;
        },
        back() {
            fetch("/api/stop?stream=" + this.streamPath);
            this.player = false;
        },
    },
};
</script>

<style scoped>
.container {
    position: relative;
    height: 500px;
    background-image: radial-gradient(#c52dd07a, #4a17987a, #0300137a);
    color: #ffffff;
    background-color: black;
    overflow: auto;
}

.search {
    padding: 10px 0;
}
</style>
