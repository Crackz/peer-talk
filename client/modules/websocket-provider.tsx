"use client";
import { connectSocket } from "@/api";
import React, { createContext, useEffect, useRef, useState } from "react";
import SimplePeer from "simple-peer";

export enum MessageType {
  BOT = "BOT_MESSAGE",
  USER = "USER_MESSAGE",
  PEER_INIT = "PEER_INIT_MESSAGE",
  PEER_SIGNAL = "PEER_SIGNAL_MESSAGE",
}

export enum InitPeerType {
  RECEIVE = "PEER_RECEIVE_INIT_TYPE",
  SEND = "PEER_SEND_INIT_TYPE",
  CLOSE = "PEER_CLOSE_INIT_TYPE",
}

export interface SocketMessage {
  id: string;
  type: MessageType;
}

export interface SocketTextMessage extends SocketMessage {
  payload: {
    text: string;
    username: string;
    createdAt: string;
  };
}

export interface SocketPeerInitMessage extends SocketMessage {
  payload: {
    userId: string;
    initType: InitPeerType;
  };
}

export interface SocketPeerSignalMessage extends SocketMessage {
  payload: {
    userId: string;
    signal: SimplePeer.SignalData;
  };
}

export interface SendSocketCreateOrJoinMessage {
  eventName: "createOrJoinRoom";
  payload: {
    roomId: string;
  };
}

export interface SendSocketTextMessage {
  eventName: "message";
  payload: {
    text: string;
  };
}
export interface SendSocketSignalMessage {
  eventName: "signal";
  payload: {
    signal: SimplePeer.SignalData;
    receiverId: string;
  };
}

export interface SendSocketPeerInitMessage {
  eventName: "peerInit";
  payload: {
    receiverId: string;
  };
}

export type SendMessage =
  | SendSocketCreateOrJoinMessage
  | SendSocketTextMessage
  | SendSocketSignalMessage
  | SendSocketPeerInitMessage;

export type ReceivedMessage =
  | SocketTextMessage
  | SocketPeerInitMessage
  | SocketPeerSignalMessage;

export const WebsocketContext = createContext<{
  lastSocketMessage: ReceivedMessage | null;
  sendSocketMessage: (msg: SendMessage) => void;
  isSocketConnected: boolean;
}>({
  lastSocketMessage: null,
  sendSocketMessage: () => {},
  isSocketConnected: false,
});

interface WebSocketProviderProps {
  children: React.ReactNode;
}

const WebSocketProvider = ({ children }: WebSocketProviderProps) => {
  const [isSocketConnected, setIsSocketConnected] = useState(false);
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [lastSocketMessage, setLastSocketMessage] =
    useState<ReceivedMessage | null>(null);

  useEffect(() => {
    connectSocket()
      .then((conn) => {
        conn.onmessage = (messageEvent: MessageEvent) => {
          console.log("GOT THIS MESSAGE: ", JSON.parse(messageEvent.data));
          setLastSocketMessage(
            JSON.parse(messageEvent.data) as ReceivedMessage
          );
        };

        setIsSocketConnected(true);
        setSocket(conn);
      })
      .catch((err) => {
        console.log("FAILED TO CONNECT : ", err);
      });

    return () => {
      if (socket) {
        socket.close();
      }
    };
  }, []);

  const sendMessage = (msg: SendMessage) => {
    if (!socket) {
      console.log(
        "Couldn't send a message because socket is not initialized: ",
        msg
      );
      return;
    }

    console.log("SENDING THIS MESSAGE: ", msg);
    socket.send(JSON.stringify(msg));
  };

  return (
    <WebsocketContext.Provider
      value={{
        isSocketConnected,
        lastSocketMessage,
        sendSocketMessage: sendMessage,
      }}
    >
      {children}
    </WebsocketContext.Provider>
  );
};

export default WebSocketProvider;
