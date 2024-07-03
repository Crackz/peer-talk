import { create } from "zustand";

type PeerNotificationsDetails = {
  shouldSendPeerInitMessage: boolean;
};

type UsersPeersNotificationsState = {
  usersPeersNotifications: Record<string, PeerNotificationsDetails>;
  setShouldSendPeerInitMessage: (
    userId: string,
    shouldSendPeerInitMessage: boolean
  ) => void;
  resetPeersNotifications: () => void;
};

export const useUsersPeersNotificationsStore =
  create<UsersPeersNotificationsState>((set) => ({
    usersPeersNotifications: {},
    setShouldSendPeerInitMessage: (
      userId: string,
      shouldSendPeerInitMessage: boolean
    ) =>
      set(
        (
          state
        ): Pick<UsersPeersNotificationsState, "usersPeersNotifications"> => {
          return {
            usersPeersNotifications: {
              ...state.usersPeersNotifications,
              [userId]: {
                ...(state.usersPeersNotifications[userId] || {}),
                shouldSendPeerInitMessage,
              },
            },
          };
        }
      ),

    resetPeersNotifications: () =>
      set(() => ({
        usersPeersNotifications: {},
      })),
  }));
