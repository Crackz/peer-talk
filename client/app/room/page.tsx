"use client";
import RoomMessages from "@/components/room-messages";
import RoomStream from "@/components/room-stream";
import WebSocketProvider from "@/modules/websocket-provider";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import styles from "./styles.module.css";

/**
 * UserMedia constraints
 */
const constraints: MediaStreamConstraints = {
  audio: {
    noiseSuppression: true,
    echoCancellation: true,
  },
  video: {
    facingMode: {
      ideal: "user",
    },
  },
};

const Page = () => {
  const router = useRouter();

  const searchParams = useSearchParams();
  const roomId = searchParams.get("id");
  const [userStream, setUserStream] = useState<MediaStream | null>(null);

  useEffect(() => {
    navigator.mediaDevices
      .getUserMedia(constraints)
      .then((stream) => {
        const ctx = new AudioContext();
        const gainNode = ctx.createGain();
        const audioDest = ctx.createMediaStreamDestination();
        const source = ctx.createMediaStreamSource(stream);

        // gainNode is set to 0.5
        gainNode.connect(audioDest);
        gainNode.gain.value = 0.5;
        source.connect(gainNode);

        setUserStream(stream);
      })
      .catch((err) => {
        console.log("FAILED TO GET USER STREAM", err);
      });
  }, []);

  useEffect(() => {
    if (!roomId) {
      router.push("/lobby");
      return;
    }
  }, [roomId, router]);

  if (!roomId) {
    //TODO: Design Proper Loading
    return <h1>Loading...</h1>;
  }

  const renderRoomContent = () => {
    if (!userStream) {
      return (
        <div id={styles.enable__stream__container}>
          <div>Please, enable Camera and Microphone</div>
        </div>
      );
    }

    return (
      <div id={styles.room__container}>
        <RoomStream roomId={roomId} userStream={userStream} />
        <RoomMessages />
      </div>
    );
  };

  return (
    <WebSocketProvider>
      <main className={styles.container}>{renderRoomContent()}</main>
    </WebSocketProvider>
  );
};

export default Page;
