<template>
    <div class="player-wrap">
        <template v-if="rtcStream">
            <video :srcObject.prop="rtcStream" autoplay muted controls></video>
        </template>
    </div>
</template>

<script>
    export default {
        name: "WebrtcPlayer",
        rtcPeerConnection: null,
        data() {
            return {
                iceConnectionState: '',
                rtcPeerConnectionInit: false,
                rtcStream: null
            }
        },
        props: {
            streamPath: {
                type: String,
                default: ''
            }
        },
        async created() {
            await this.initRtcPeerConnection();
            console.log('initRtcPeerConnectioned');
            if (this.streamPath) {
                await this.play(this.streamPath);
                console.log('played');
            }
        },
        methods: {
            async initRtcPeerConnection() {
                const rtcPeerConnection = new RTCPeerConnection();

                rtcPeerConnection.addTransceiver('video', {
                    direction: "recvonly"
                });

                rtcPeerConnection.onsignalingstatechange = e => {
                    console.log('onsignalingstatechange', e);
                };

                rtcPeerConnection.oniceconnectionstatechange = e => {
                    console.log('oniceconnectionstatechange', rtcPeerConnection.iceConnectionState);
                };

                rtcPeerConnection.onicecandidate = event => {
                    console.log('onicecandidate', event);
                };

                rtcPeerConnection.ontrack = event => {
                    console.log('ontrack', event);
                    if (event.track.kind === "video") {
                        this.rtcStream = event.streams[0];
                    }
                };

                const rtcSessionDescriptionInit = await rtcPeerConnection.createOffer();
                await rtcPeerConnection.setLocalDescription(rtcSessionDescriptionInit);
                this.rtcPeerConnectionInit = true;
                this.$options.rtcPeerConnection = rtcPeerConnection;
            },

            //
            async play(streamPath) {
                const rtcPeerConnection = this.$options.rtcPeerConnection;
                const localDescriptionData = rtcPeerConnection.localDescription.toJSON();
                const result = await this.ajax({
                    type: "POST",
                    processData: false,
                    data: JSON.stringify(localDescriptionData),
                    url: "/webrtc/play?streamPath=" + streamPath,
                    dataType: "json"
                });
                if (result.errmsg) {
                    console.error(result.errmsg);
                    return;
                }
                //
                rtcPeerConnection.setRemoteDescription(new RTCSessionDescription({
                    type: result.type,
                    sdp: result.sdp
                }));
            },
            close() {
                const rtcPeerConnection = this.$options.rtcPeerConnection;
                rtcPeerConnection && rtcPeerConnection.close();
            }
        },
        destroyed() {
            this.close();
        }
    }
</script>

<style scoped>
    .player-wrap {
        width: 100%;
        height: 100%;
        border-radius: 4px;
        box-shadow: 0 0 5px #40d3fc, inset 0 0 5px #40d3fc, 0 0 0 1px #40d3fc;
    }

    .player-wrap video {
        width: 100%;
        height: 100%;
    }
</style>
