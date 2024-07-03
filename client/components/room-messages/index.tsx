"use client";
import {
  SocketTextMessage,
  MessageType,
  WebsocketContext,
  SendSocketTextMessage,
} from "@/modules/websocket-provider";
import { useContext, useEffect, useRef, useState } from "react";
import styles from "./styles.module.css";

interface RoomMessagesProps {}

const RoomMessages = ({}: RoomMessagesProps) => {
  const { lastSocketMessage, sendSocketMessage } = useContext(WebsocketContext);
  const [messages, setMessages] = useState<SocketTextMessage[]>([]);
  const messagesRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const isTextMessage =
      lastSocketMessage &&
      (lastSocketMessage.type === MessageType.BOT ||
        lastSocketMessage.type === MessageType.USER);

    if (isTextMessage) {
      setMessages([...messages, lastSocketMessage as SocketTextMessage]);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lastSocketMessage]);

  useEffect(() => {
    messagesRef.current?.lastElementChild?.scrollIntoView({
      behavior: "smooth",
    });
  }, [messages]);

  const getBotMessage = (message: SocketTextMessage) => {
    return (
      <div className={styles.message__body__bot}>
        <strong className={styles.message__author__bot}>
          🤖 {message.payload.username}
        </strong>
        <p className={styles.message__text__bot}>{message.payload.text}</p>
      </div>
    );
  };

  const getUserMessage = (message: SocketTextMessage) => {
    return (
      <div className={styles.message__body}>
        <strong className={styles.message__author}>
          {message.payload.username}
        </strong>
        <p className={styles.message__text}>{message.payload.text}</p>
      </div>
    );
  };
  return (
    <section id={styles.messages__container}>
      <div id={styles.messages} ref={messagesRef}>
        {messages.map((message) => {
          return (
            <div key={message.id} className={styles.message__wrapper}>
              {message.type === MessageType.USER
                ? getUserMessage(message)
                : getBotMessage(message)}
            </div>
          );
        })}
      </div>
      <form id={styles.message__form}>
        <input
          type="text"
          name="message"
          placeholder="Send a message...."
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();
              if (e.currentTarget.value !== "") {
                const msg: SendSocketTextMessage = {
                  eventName: "message",
                  payload: {
                    text: e.currentTarget.value,
                  },
                };
                sendSocketMessage(msg);
                e.currentTarget.value = "";
              }
            }
          }}
        />
      </form>
    </section>
  );
};

export default RoomMessages;
