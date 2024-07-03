"use client";
import Image from "next/image";
import { useUsersPeersStore } from "../../stores/users-peers.store";
import { StreamTracksStatesState } from "../room-stream";
import RoomVideo from "../room-video";
import styles from "./styles.module.css";

interface RoomStreamVideosProps {
  userStream: MediaStream | null;
  userStreamTracksStates: StreamTracksStatesState;
}

const RoomStreamVideos = ({
  userStream,
  userStreamTracksStates,
}: RoomStreamVideosProps) => {
  const usersPeers = useUsersPeersStore((state) => state.usersPeers);

  const renderMicOffIcon = () => {
    return (
      <div
        style={{
          position: "absolute",
          alignSelf: "end",
          justifySelf: "end",
          margin: "5px",
          width: "30px",
          height: "30px",
        }}
      >
        <Image fill src={`/icons/mic-off-red.svg`} alt="mic-icon" />
      </div>
    );
  };
  const renderLocalVideo = () => {
    return (
      <div className={styles.video__container}>
        <RoomVideo
          isVideoEnabled={userStreamTracksStates.camera}
          srcObject={userStream}
          autoPlay
        />
        {!userStreamTracksStates.mic && renderMicOffIcon()}
      </div>
    );
  };
  const renderRemoteVideos = () => {
    return Object.keys(usersPeers).map((userId) => {
      const { stream } = usersPeers[userId];
      return (
        <div key={userId} className={styles.video__container}>
          <RoomVideo
            isVideoEnabled={usersPeers[userId].isCameraEnabled}
            srcObject={stream || null}
            autoPlay
          />
          {!usersPeers[userId].isMicEnabled && renderMicOffIcon()}
        </div>
      );
    });
  };

  return (
    <div id={styles.streams__container}>
      {renderLocalVideo()}
      {renderRemoteVideos()}
    </div>
  );
};

export default RoomStreamVideos;
