"use client";
import styles from "./styles.module.css";

interface RoomParticipantsProps {}

const RoomParticipants = ({}: RoomParticipantsProps) => {
  return (
    <section id={styles.members__container}>
      <div id={styles.members__header}>
        <p>Participants</p>
        <strong id={styles.members__count}>0</strong>
      </div>

      <div id={styles.member__list}></div>
    </section>
  );
};

export default RoomParticipants;
