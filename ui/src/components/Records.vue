<template>
    <Modal v-bind="$attrs"
           draggable
           width="900"
           v-on="$listeners"
           title="Records"
           @on-ok="onClosePreview"
           @on-cancel="onClosePreview">
        <div class="container">
            <div class="search">
                <DatePicker type="date" :options="timeOptions" :value="search.time" placeholder="请选择时间"
                            style="width: 200px" :clearable="false"
                            @on-change="handleTimeChange"></DatePicker>
            </div>
            <div>
                <mu-data-table :columns="columns" :data="recordList">
                    <template #expand="prop">
                        <div>
                            <m-button @click="play(prop.row)">播放</m-button>
                        </div>
                    </template>
                    <template #default="scope">
                        <td>{{scope.row.DeviceID}}</td>
                        <td>{{scope.row.Name}}</td>
                        <td>{{scope.row.time}}</td>
                        <td>{{scope.row.length}}</td>
                        <td>{{scope.row.FilePath}}</td>
                        <td>{{scope.row.Address}}</td>
                        <td>{{scope.row.Type}}</td>
                    </template>
                </mu-data-table>
            </div>
        </div>
    </Modal>
</template>

<script>
    import {getOneTimeRange, formatTimeTips, parseTime} from "../utils";

    const _now = new Date();

    export default {
        name: "Records",
        props: {
            records: Array
        },
        data() {
            return {
                timeOptions: {
                    disabledDate(date) {
                        return date && date.valueOf() > Date.now();
                    }
                },
                search: {
                    id: null,
                    channel: 0,
                    deviceId: null,
                    time: _now
                },
                columns: Object.freeze(
                    [
                        '设备ID',
                        '名称',
                        '时间',
                        '时长',
                        '文件路径',
                        '地址',
                        '类型'
                    ].map((title) => ({
                        title,
                    }))
                )
            }
        },
        computed: {
            startTime() {
                if (!this.search.time) {
                    return ''
                }
                const start = getOneTimeRange(this.search.time).start;
                return new Date(start).toISOString();
            },
            endTime() {

                if (!this.search.time) {
                    return ''
                }

                const end = getOneTimeRange(this.search.time).end;
                return new Date(end).toISOString();
            },
            recordList() {
                const list = this.records.map((record) => {
                    const startTime = new Date(record.StartTime).getTime();
                    const endTime = new Date(record.EndTime).getTime();
                    const timestamp = endTime - startTime;
                    const timeLength = formatTimeTips(timestamp);

                    const _startTime = parseTime(startTime);

                    record.length = timeLength;
                    record.time = _startTime;
                    return record;
                });
                return list;
            }
        },

        methods: {
            getList(options) {
                this.search.id = options.id;
                this.search.channel = options.channel;
                this.search.deviceId = options.deviceId;
                this._fetchList();
            },

            _fetchList() {
                if (this.search.id && this.search.channel && this.startTime && this.endTime) {
                    const query = {
                        id: this.search.id,
                        channel: this.search.channel,
                        startTime: this.startTime,
                        endTime: this.endTime
                    };

                    this.ajax.get("/gb28181/query/records", query).then((x) => {

                    });
                }
            },
            onClosePreview() {
                this.$emit('close');
            },
            handleTimeChange(date) {
                this.search.time = new Date(date);
                this._fetchList();
            },
            play() {

            }
        }
    }
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
