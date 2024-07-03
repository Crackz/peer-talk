import { StreamTracksStatesState } from "@/components/room-stream";
import { create } from "zustand";

type PeerNotificationsDetails = {
  camera: boolean;
  mic: boolean;
};

export type UsersPeersNotificationsState = {
  usersPeersStreamTracks: Record<string, PeerNotificationsDetails>;
  setStreamTrackState: (
    userId: string,
    streamTrackState: Partial<StreamTracksStatesState>
  ) => void;
};

export const useUsersPeersStreamsTracksStore =
  create<UsersPeersNotificationsState>((set) => ({
    usersPeersStreamTracks: {},
    setStreamTrackState: (
      userId: string,
      streamTrackState: Partial<StreamTracksStatesState>
    ) =>
      set(
        (
          state
        ): Pick<UsersPeersNotificationsState, "usersPeersStreamTracks"> => {
          return {
            usersPeersStreamTracks: {
              ...state.usersPeersStreamTracks,
              [userId]: {
                ...(state.usersPeersStreamTracks[userId] || {}),
                ...streamTrackState,
              },
            },
          };
        }
      ),
  }));
