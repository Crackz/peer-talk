"use client";
import CreateOrJoinLobbyRoom from "@/components/lobby-room";
import { useRouter } from "next/navigation";

const Page = () => {
  const router = useRouter();

  const submitHandler = (roomId: string) => {
    router.push(`/room?id=${roomId}`);
  };

  return <CreateOrJoinLobbyRoom onSubmit={submitHandler} />;
};

export default Page;
