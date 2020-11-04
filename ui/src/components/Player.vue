<template>
  <Modal
    v-bind="$attrs"
    draggable
    width="750"
    v-on="$listeners"
    :title="streamPath"
    @on-ok="onClosePreview"
    @on-cancel="onClosePreview"
  >
    <div class="container">
      <video ref="webrtc" :srcObject.prop="stream" width="488" height="275" autoplay muted controls></video>
      <div class="control">
        <svg @click="$emit('ptz',n)" v-for="n in 8" :class="'arrow'+n" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="64" height="64"><defs><style type="text/css"></style></defs><path d="M682.666667 955.733333H341.333333a17.066667 17.066667 0 0 1-17.066666-17.066666V529.066667H85.333333a17.066667 17.066667 0 0 1-12.066133-29.1328l426.666667-426.666667a17.0496 17.0496 0 0 1 24.132266 0l426.666667 426.666667A17.066667 17.066667 0 0 1 938.666667 529.066667H699.733333v409.6a17.066667 17.066667 0 0 1-17.066666 17.066666z m-324.266667-34.133333h307.2V512a17.066667 17.066667 0 0 1 17.066667-17.066667h214.801066L512 109.4656 126.532267 494.933333H341.333333a17.066667 17.066667 0 0 1 17.066667 17.066667v409.6z" p-id="6849"></path></svg>
      </div>
    </div>
    <div slot="footer">
      <mu-badge v-if="remoteSDP">
        <a slot="content" :href="remoteSDPURL" download="remoteSDP.txt">remoteSDP</a>
      </mu-badge>
      <mu-badge v-if="localSDP">
        <a slot="content" :href="localSDPURL" download="localSDP.txt">localSDP</a>
      </mu-badge>
    </div>
  </Modal>
</template>
<script>
let pc = null;
export default {
  data() {
    return {
      iceConnectionState: pc && pc.iceConnectionState,
      stream: null,
      localSDP: "",
      remoteSDP: "",
      remoteSDPURL: "",
      localSDPURL: "",
      streamPath: ""
    };
  },
props:{
  PublicIP:String
},
    methods: {
        async play(streamPath) {
            pc = new RTCPeerConnection();
            pc.addTransceiver('video',{
              direction:'recvonly'
            })
            this.streamPath = streamPath;
            pc.onsignalingstatechange = e => {
                //console.log(e);
            };
            pc.oniceconnectionstatechange = e => {
                this.$toast.info(pc.iceConnectionState);
                this.iceConnectionState = pc.iceConnectionState;
            };
            pc.onicecandidate = event => {
                console.log(event)
            };
            pc.ontrack = event => {
               // console.log(event);
                if (event.track.kind == "video")
                    this.stream = event.streams[0];
            };
            await pc.setLocalDescription(await pc.createOffer());
            this.localSDP = pc.localDescription.sdp;
            this.localSDPURL = URL.createObjectURL(
                new Blob([this.localSDP], { type: "text/plain" })
            );
            const result = await this.ajax({
                type: "POST",
                processData: false,
                data: JSON.stringify(pc.localDescription.toJSON()),
                url: "/webrtc/play?streamPath=" + this.streamPath,
                dataType: "json"
            });
            if (result.errmsg) {
                this.$toast.error(result.errmsg);
                return;
            } else {
                this.remoteSDP = result.sdp;
                this.remoteSDPURL = URL.createObjectURL(new Blob([this.remoteSDP], { type: "text/plain" }));
            }
            await pc.setRemoteDescription(new RTCSessionDescription(result));
        },
        onClosePreview() {
            pc.close();
        }
    }
};
</script>
<style scoped>
  .arrow1{
    grid-column: 2;
    grid-row: 1;
  }
  .arrow2{
    transform: rotate(90deg);
    grid-column: 3;
    grid-row: 2;
  }
  .arrow3{
    transform: rotate(180deg);
    grid-column: 2;
    grid-row: 3;
  }
  .arrow4{
    transform: rotate(270deg);
    grid-column: 1;
    grid-row: 2;
  }

  .arrow5{
    transform: rotate(-45deg);
    grid-column: 1;
    grid-row: 1;
  }

  .arrow6{
    transform: rotate(45deg);
    grid-column: 3;
    grid-row: 1;
  }

  .arrow7{
    transform: rotate(-135deg);
    grid-column: 1;
    grid-row: 3;
  }

  .arrow8{
    transform: rotate(135deg);
    grid-column: 3;
    grid-row: 3;
  }

.container {
  position: relative;
}
.control {
  position: absolute;
  top: 20px;
  right: 0;
  display: grid;
  grid-template-columns: repeat(3, 33.33%);
  grid-template-rows: repeat(3, 33.33%);
  width: 192px;
  height: 192px;
}
.control >* {
  cursor: pointer;
  fill: gray;
  width: 50px;
  height: 50px;
}
.control >*:hover{
  fill: cyan
}
</style>
