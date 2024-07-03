export const getRandomNumber = () => Math.floor(Math.random() * 100000);

const connectSocket = async (): Promise<WebSocket> => {
  const accessToken = localStorage.getItem("accessToken");
  const socket = new WebSocket(
    `${process.env.NEXT_PUBLIC_WEBSOCKET_URL}/ws?authorization="Bearer ${accessToken}`
  );
  console.log("Connecting...");

  return await new Promise((resolve, reject) => {
    socket.onopen = () => {
      console.log("Successfully Connected");
      resolve(socket);
    };

    socket.onclose = (e) => {
      console.log("SOCKET CLOSED");
      reject(e);
    };

    socket.onerror = (e) => {
      reject(e);
    };
  });
};

export { connectSocket };
