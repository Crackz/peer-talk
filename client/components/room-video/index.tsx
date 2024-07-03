import Image from "next/image";
import { VideoHTMLAttributes, useEffect, useRef } from "react";

type PropsType = VideoHTMLAttributes<HTMLVideoElement> & {
  srcObject: MediaStream | null;
  isVideoEnabled: boolean;
};

export default function RoomVideo({
  isVideoEnabled,
  srcObject,
  ...props
}: PropsType) {
  const refVideo = useRef<HTMLVideoElement>(null);

  const openPictureMode = () => {
    if (refVideo.current) {
      refVideo.current.requestPictureInPicture();
    }
  };

  useEffect(() => {
    if (!refVideo.current || !srcObject) return;
    refVideo.current.srcObject = srcObject;
  }, [srcObject, isVideoEnabled]);

  const renderContent = () => {
    if (!srcObject || !isVideoEnabled) {
      return (
        <Image
          src={"/icons/camera-off.svg"}
          alt="camera-icon"
          height={100}
          width={100}
        />
      );
    }

    return (
      <video
        style={{
          flex: 1,
          objectFit: "contain",
          cursor: "pointer",
        }}
        ref={refVideo}
        onClick={() => openPictureMode()}
        onTouchStart={() => openPictureMode()}
        {...props}
      />
    );
  };

  return (
    <div
      style={{
        display: "flex",
        flex: 1,
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      {renderContent()}
    </div>
  );
}
