"use client";
import {
  InitPeerType,
  MessageType,
  SendSocketCreateOrJoinMessage,
  SendSocketPeerInitMessage,
  SendSocketSignalMessage,
  SocketPeerInitMessage,
  SocketPeerSignalMessage,
  WebsocketContext,
} from "@/modules/websocket-provider";
import { useContext, useEffect, useState } from "react";
import SimplePeer from "simple-peer";
import { useUsersPeersNotificationsStore } from "../../stores/users-peers-notification.store";
import { useUsersPeersStore } from "../../stores/users-peers.store";
import RoomStreamButtons from "../room-stream-buttons";
import RoomStreamVideos from "../room-stream-videos";
import styles from "./styles.module.css";

interface RoomStreamProps {
  roomId: string;
  userStream: MediaStream;
}

export type StreamTracksStatesState = {
  mic: boolean;
  camera: boolean;
};

const configuration = {
  // Using From https://www.metered.ca/tools/openrelay/
  iceServers: [
    {
      urls: "stun:openrelay.metered.ca:80",
    },
    {
      urls: "turn:openrelay.metered.ca:80",
      username: "openrelayproject",
      credential: "openrelayproject",
    },
    {
      urls: "turn:openrelay.metered.ca:443",
      username: "openrelayproject",
      credential: "openrelayproject",
    },
    {
      urls: "turn:openrelay.metered.ca:443?transport=tcp",
      username: "openrelayproject",
      credential: "openrelayproject",
    },
  ],
};

const RoomStream = ({ roomId, userStream }: RoomStreamProps) => {
  const { isSocketConnected, sendSocketMessage, lastSocketMessage } =
    useContext(WebsocketContext);

  const {
    usersPeers,
    addPeer: addPeerToUsersPeers,
    removePeer: removePeerFromUsersPeers,
    addStream: addStreamToUsersPeers,
    removeAllPeers: removeAllPeersFromUsersPeers,
    updateStreamTrackState: updateStreamTrackStateToUsersPeers,
  } = useUsersPeersStore();

  const {
    usersPeersNotifications,
    setShouldSendPeerInitMessage,
    resetPeersNotifications,
  } = useUsersPeersNotificationsStore();

  const [streamTracksStates, setStreamTracksStates] =
    useState<StreamTracksStatesState>({
      camera: userStream.getVideoTracks().some((track) => track.enabled),
      mic: userStream.getAudioTracks().some((track) => track.enabled),
    });

  const toggleStreamTrackState = (type: keyof StreamTracksStatesState) => {
    let tracks: MediaStreamTrack[] = [];
    switch (type) {
      case "camera":
        tracks = userStream.getVideoTracks();
        break;
      case "mic":
        tracks = userStream.getAudioTracks();
        break;
      default:
        throw new Error("Unknown track type: " + type);
    }

    for (const track of tracks) {
      track.enabled = !track.enabled;
    }

    Object.keys(usersPeers).forEach((userId) => {
      const msg: PeerStreamTrackState = {
        camera: userStream.getVideoTracks().some((track) => track.enabled),
        mic: userStream.getAudioTracks().some((track) => track.enabled),
      };
      usersPeers[userId].peer.send(JSON.stringify(msg));
    });

    setStreamTracksStates({
      ...streamTracksStates,
      [type]: tracks.some((track) => track.enabled),
    });
  };

  const createPeer = (userId: string, isInitializer: boolean): void => {
    console.log("PEER USER ID: ", userId, isInitializer);
    const peer = new SimplePeer({
      initiator: isInitializer,
      stream: userStream!,
      config: configuration,
    });

    peer.on("signal", (data: SimplePeer.SignalData) => {
      const msg: SendSocketSignalMessage = {
        eventName: "signal",
        payload: {
          signal: data,
          receiverId: userId,
        },
      };
      sendSocketMessage(msg);
    });

    peer.on("stream", (stream) => {
      console.log(`RECEIVED STREAM FROM USER ID ${userId}`);

      console.log(
        "IS ENABLED VIDEO AND MIC: ",
        stream.getVideoTracks().some((track) => track.enabled),
        stream.getAudioTracks().some((track) => track.enabled)
      );
      addStreamToUsersPeers(userId, stream);
    });

    peer.on("data", (data: string) => {
      const streamTrackState: PeerStreamTrackState = JSON.parse(data);
      console.log("streamTrackState: ", streamTrackState);
      updateStreamTrackStateToUsersPeers(userId, streamTrackState);
    });

    addPeerToUsersPeers(userId, peer);

    if (!isInitializer) {
      setShouldSendPeerInitMessage(userId, true);
    }
  };

  const handlePeerInit = (payload: SocketPeerInitMessage["payload"]) => {
    switch (payload.initType) {
      case InitPeerType.RECEIVE: {
        console.log("SEND PEER INIT TO : ", payload.userId);
        createPeer(payload.userId, false);
        break;
      }
      case InitPeerType.SEND: {
        console.log("RECEIVED PEER INIT FROM : ", payload.userId);
        createPeer(payload.userId, true);
        break;
      }

      case InitPeerType.CLOSE: {
        console.log("REMOVING PEER: ", payload.userId);

        removePeerFromUsersPeers(payload.userId);
        break;
      }
      default:
        throw new Error("Unhandled init type: " + payload.initType);
    }
  };

  const handlePeerSignal = (payload: SocketPeerSignalMessage["payload"]) => {
    usersPeers[payload.userId].peer.signal(payload.signal);
  };

  useEffect(() => {
    if (!lastSocketMessage) {
      return;
    }

    switch (lastSocketMessage.type) {
      case MessageType.PEER_INIT: {
        const { payload } = lastSocketMessage as SocketPeerInitMessage;
        handlePeerInit(payload);
        break;
      }
      case MessageType.PEER_SIGNAL: {
        const { payload } = lastSocketMessage as SocketPeerSignalMessage;
        handlePeerSignal(payload);
        break;
      }
    }
  }, [lastSocketMessage]);

  useEffect(() => {
    if (isSocketConnected && userStream) {
      console.log("LISTENING TO SOCKET MESSAGES FOR ROOM:  ", roomId);

      const createOrJoinRoomPayload: SendSocketCreateOrJoinMessage = {
        eventName: "createOrJoinRoom",
        payload: {
          roomId,
        },
      };
      sendSocketMessage(createOrJoinRoomPayload);
    }

    return () => {
      Object.keys(usersPeers).forEach((userId) => {
        console.log("DESTROYING PEER: ", userId);
        usersPeers[userId].peer.destroy();
      });

      removeAllPeersFromUsersPeers();
      resetPeersNotifications();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isSocketConnected, userStream]);

  useEffect(() => {
    Object.keys(usersPeersNotifications).forEach((userId) => {
      if (usersPeersNotifications[userId].shouldSendPeerInitMessage) {
        console.log("Sending Peer Init Message to user: ", userId);

        const msg: SendSocketPeerInitMessage = {
          eventName: "peerInit",
          payload: {
            receiverId: userId,
          },
        };
        sendSocketMessage(msg);
        setShouldSendPeerInitMessage(userId, false);
      }
    });
  }, [usersPeersNotifications]);

  if (!isSocketConnected) {
    return <p>Connecting...</p>;
  }

  return (
    <section id={styles.stream__container}>
      <div id={styles.stream__box}></div>
      <RoomStreamVideos
        userStream={userStream}
        userStreamTracksStates={streamTracksStates}
      />
      <RoomStreamButtons
        userStream={userStream}
        toggleStreamTrackState={toggleStreamTrackState}
        streamTracksStates={streamTracksStates}
      />
    </section>
  );
};

export default RoomStream;
