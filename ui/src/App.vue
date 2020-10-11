<template>
    <div>
        <mu-data-table :data="Devices" :columns="columns">
            <template #expand="prop" v-if="prop.row.Channels">
                <mu-data-table :data="prop.row.Channels" :columns="columns2">
                    <template #default="{ row: item, $index }">
                        <td>{{ item.DeviceID }}</td>
                        <td>{{ item.Name }}</td>
                        <td>{{ item.Manufacturer }}</td>
                        <td>{{ item.Address }}</td>
                        <td>{{ item.Status }}</td>
                        <td>
                            <mu-button flat @click="invite(prop.row.ID, $index)"
                                >invite</mu-button
                            >
                            <!--          <mu-button flat @click="preview(item)">预览</mu-button>-->
                            <mu-button flat @click="left(prop.row.ID, $index)"
                                >👈</mu-button
                            >
                            <mu-button flat @click="right(prop.row.ID, $index)"
                                >👉</mu-button
                            >
                        </td>
                    </template>
                </mu-data-table>
            </template>
            <template #default="{ row: item }">
                <td>{{ item.ID }}</td>
                <td>{{ item.Channels?item.Channels.length:0}}</td>
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
</template>
<script>
export default {
    data() {
        return {
            Devices: [],
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
            ]).map((title) => ({ title })),
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
        preview(item) {
            alert("尚未实现");
        },
        left(id, channel) {
            this.ajax.get("/gb28181/control", {
                id,
                channel,
                ptzcmd: "A50F01021F0000D6",
            });
        },
        right(id, channel) {
            this.ajax.get("/gb28181/control", {
                id,
                channel,
                ptzcmd: "A50F01011F0000D6",
            });
        },
        invite(id, channel) {
            this.ajax.get("/gb28181/invite", { id, channel });
        },
    },
};
</script>