"use client";
import Image from "next/image";
import styles from "./styles.module.css";
import { StreamTracksStatesState } from "../room-stream";

interface RoomStreamButtonsProps {
  userStream: MediaStream | null;
  streamTracksStates: StreamTracksStatesState;
  toggleStreamTrackState: (type: keyof StreamTracksStatesState) => void;
}
const RoomStreamButtons = ({
  userStream,
  streamTracksStates,
  toggleStreamTrackState,
}: RoomStreamButtonsProps) => {
  if (!userStream) return null;

  return (
    <div className={styles.stream__actions}>
      <button
        className={streamTracksStates.camera ? styles.active : ""}
        onClick={() => toggleStreamTrackState("camera")}
      >
        <Image
          src={`/icons/${
            streamTracksStates.camera ? "camera-on.svg" : "camera-off.svg"
          }`}
          alt="camera-icon"
          height={30}
          width={30}
        />
      </button>
      <button
        className={streamTracksStates.mic ? styles.active : ""}
        onClick={() => toggleStreamTrackState("mic")}
      >
        <Image
          src={`/icons/${
            streamTracksStates.mic ? "mic-on.svg" : "mic-off.svg"
          }`}
          alt="mic-icon"
          height={30}
          width={30}
        />
      </button>
    </div>
  );
};

export default RoomStreamButtons;
