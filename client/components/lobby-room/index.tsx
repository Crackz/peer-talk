"use client";
import { ChangeEvent, useState } from "react";
import styles from "./styles.module.css";

interface CreateOrJoinLobbyRoomProps {
  onSubmit: (roomId: string) => void;
}

const CreateOrJoinLobbyRoom = ({ onSubmit }: CreateOrJoinLobbyRoomProps) => {
  const [roomId, setRoomId] = useState("");

  const handleRoomInput = (e: ChangeEvent<HTMLInputElement>) => {
    const fieldValue = e.target.value;
    setRoomId(fieldValue);
  };

  return (
    <div id={styles.form__container}>
      <div id={styles.form__container__header}>
        <p>👋 Create or Join Room</p>
      </div>

      <form
        id={styles.lobby__form}
        onSubmit={(e) => {
          e.preventDefault();
          onSubmit(roomId);
        }}
      >
        <div className={styles.form__field__wrapper}>
          <label>Room Name</label>
          <input
            type="text"
            name="room"
            placeholder="Enter room name..."
            onChange={handleRoomInput}
          />
        </div>

        <div className={styles.form__field__wrapper}>
          <button type="submit">
            Go to Room
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
            >
              <path d="M13.025 1l-2.847 2.828 6.176 6.176h-16.354v3.992h16.354l-6.176 6.176 2.847 2.828 10.975-11z" />
            </svg>
          </button>
        </div>
      </form>
    </div>
  );
};

export default CreateOrJoinLobbyRoom;
