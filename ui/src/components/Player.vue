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
      <div class="control control2">
        <svg @click="$emit('ptz',9)" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M994.990643 859.352971L713.884166 578.246494A381.208198 381.208198 0 0 0 767.307984 383.653992C767.307984 171.765089 595.542895 0 383.653992 0S0 171.765089 0 383.653992s171.765089 383.653992 383.653992 383.653992c71.119859 0 137.507985-19.694238 194.592502-53.423818l281.106477 281.090491a95.913498 95.913498 0 1 0 135.637672-135.621686zM383.653992 671.394486c-158.912681 0-287.740494-128.827813-287.740494-287.740494S224.741311 95.913498 383.653992 95.913498s287.740494 128.827813 287.740494 287.740494-128.827813 287.740494-287.740494 287.740494z m159.85583-335.697243h-111.899081v-111.899081a47.956749 47.956749 0 1 0-95.913498 0v111.899081h-111.899081a47.956749 47.956749 0 1 0 0 95.913498h111.899081v111.899081a47.956749 47.956749 0 1 0 95.913498 0v-111.899081h111.899081a47.956749 47.956749 0 1 0 0-95.913498z"  /></svg>
        <svg @click="$emit('ptz',10)" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M994.990643 859.352971L713.884166 578.246494A381.208198 381.208198 0 0 0 767.307984 383.653992C767.307984 171.765089 595.542895 0 383.653992 0S0 171.765089 0 383.653992s171.765089 383.653992 383.653992 383.653992c71.119859 0 137.507985-19.694238 194.592502-53.423818l281.106477 281.090491a95.913498 95.913498 0 1 0 135.637672-135.621686zM383.653992 671.394486c-158.912681 0-287.740494-128.827813-287.740494-287.740494S224.741311 95.913498 383.653992 95.913498s287.740494 128.827813 287.740494 287.740494-128.827813 287.740494-287.740494 287.740494z m159.85583-335.697243H223.798162a47.956749 47.956749 0 1 0 0 95.913498h319.71166a47.956749 47.956749 0 1 0 0-95.913498z"  /></svg>
      </div>
      <div class="control control3">
        <svg @click="$emit('ptz',11)" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M956.39 400.827C922.164 266.675 828.186 155.703 701.502 99.874l94.522 443.782L956.39 400.827zM206.208 189.167C106.183 286.191 56.845 424.181 72.696 562.659l351.347-309.096-217.835-64.396zM643.118 78.847a446.363 446.363 0 0 0-138.947-16.775 448.047 448.047 0 0 0-250.583 86.934l437.868 146.949-48.338-217.108zM83.786 623.979c34.443 133.772 128.248 244.407 254.583 300.291l-95.915-426.55L83.786 623.979zM969.893 496.089a372.746 372.746 0 0 0-2.37-34.138l-329.972 303.78 196.157 69.256c91.522-88.456 141.056-211.704 136.185-338.898zM396.862 945.166a447.857 447.857 0 0 0 139.077 16.766 447.784 447.784 0 0 0 250.322-86.718L349.286 733.05l47.576 212.116z"  /><path fill="#333333" d="M397.253 471.171h245.668c22.593 0 40.923 18.32 40.923 40.913 0 22.592-18.33 40.922-40.923 40.922H397.253c-22.592 0-40.922-18.33-40.922-40.922 0-22.593 18.33-40.913 40.922-40.913z"  /><path fill="#333333" d="M479.17 634.879V389.21c0-22.593 18.32-40.923 40.913-40.923s40.923 18.33 40.923 40.923v245.668c0 22.592-18.33 40.922-40.923 40.922s-40.913-18.329-40.913-40.921z"  /></svg>
        <svg @click="$emit('ptz',12)" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M956.39 400.827C922.164 266.675 828.186 155.703 701.502 99.874l94.522 443.782L956.39 400.827z m-750.182-211.66C106.183 286.191 56.845 424.181 72.696 562.659l351.347-309.096-217.835-64.396z m436.91-110.32a446.363 446.363 0 0 0-138.947-16.775 448.047 448.047 0 0 0-250.583 86.934l437.868 146.949-48.338-217.108zM83.786 623.979c34.443 133.772 128.248 244.407 254.583 300.291l-95.915-426.55L83.786 623.979z m886.107-127.89a372.746 372.746 0 0 0-2.37-34.138l-329.972 303.78 196.157 69.256c91.522-88.456 141.056-211.704 136.185-338.898zM396.862 945.166a447.857 447.857 0 0 0 139.077 16.766 447.784 447.784 0 0 0 250.322-86.718L349.286 733.05l47.576 212.116z m0.391-474.039h245.668c22.593 0 40.923 18.32 40.923 40.912 0 22.593-18.33 40.923-40.923 40.923H397.253c-22.592 0-40.922-18.33-40.922-40.923 0-22.592 18.33-40.912 40.922-40.912z"  /></svg>
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

.control2{
  top: 210px;
}

.control3{
  top: 260px;
}

.control >* {
  cursor: pointer;
  fill: gray;
  width: 50px;
  height: 50px;
}

.control2 >*,.control3 >*{
  width: 40px;
  height: 40px;
}

.control >*:hover{
  fill: cyan
}
</style>
