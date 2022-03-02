let stream = new MediaStream();
let config = {
  iceServers: [
    {
      urls: ["stun:stun.l.google.com:19302"],
    },
  ],
};

const pc = new RTCPeerConnection(config);
// 媒体协商时触发事件
pc.onnegotiationneeded = async () => {
  const offer = await pc.createOffer();
  await pc.setLocalDescription(offer);
  getRemoteSdp();
};

let log = (msg) => {
  document.getElementById("div").innerHTML += msg + "<br>";
};

// 接收到流媒体
pc.ontrack = function (event) {
  stream.addTrack(event.track);
  videoElem.srcObject = stream;
  log(event.streams.length + " track is delivered");
};
pc.oniceconnectionstatechange = (e) => log(pc.iceConnectionState);

$(document).ready(function () {
  $("#submit").click((e) => {
    console.log("submit");
    // 只建立轨道，不调用本地摄像头麦克风
    pc.addTransceiver("video", {
      direction: "sendrecv",
    });
  });
});

function getRemoteSdp() {
  const rtspUrl = $("#rtspUrl").val();
  $.post(
    "../receiver",
    {
      rtspUrl,
      data: btoa(pc.localDescription.sdp),
    },
    function (sdp) {
      try {
        pc.setRemoteDescription(
          new RTCSessionDescription({ type: "answer", sdp: atob(sdp) })
        );
      } catch (e) {
        console.warn(e);
      }
    }
  );
}
