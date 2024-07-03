import SimplePeer from "simple-peer";
import { create } from "zustand";

type PeerDetails = {
  peer: SimplePeer.Instance;
  isCameraEnabled: boolean;
  isMicEnabled: boolean;
  stream?: MediaStream;
};

export type UsersPeersState = {
  usersPeers: Record<string, PeerDetails>;
  addPeer: (userId: string, peer: SimplePeer.Instance) => void;
  removePeer: (userId: string) => void;
  addStream: (userId: string, stream: MediaStream) => void;
  removeAllPeers: () => void;
  updateStreamTrackState: (
    userId: string,
    streamTrackState: PeerStreamTrackState
  ) => void;
};

export const useUsersPeersStore = create<UsersPeersState>((set) => ({
  usersPeers: {},
  addPeer: (userId: string, peer: SimplePeer.Instance) =>
    set((state) => ({
      usersPeers: {
        ...state.usersPeers,
        [userId]: {
          ...state.usersPeers[userId],
          isCameraEnabled: false,
          isMicEnabled: false,
          peer,
        },
      },
    })),
  removePeer: (userId: string) =>
    set((state) => {
      const usersPeersCopy = { ...state.usersPeers };
      delete usersPeersCopy[userId];

      return { usersPeers: usersPeersCopy };
    }),

  removeAllPeers: () => set(() => ({ usersPeers: {} })),

  addStream: (userId: string, stream: MediaStream) =>
    set((state) => ({
      usersPeers: {
        ...state.usersPeers,
        [userId]: {
          ...state.usersPeers[userId],
          isCameraEnabled: stream
            .getVideoTracks()
            .some((track) => track.enabled),
          isMicEnabled: stream.getAudioTracks().some((track) => track.enabled),
          stream,
        },
      },
    })),

  updateStreamTrackState: (
    userId: string,
    streamTrackState: PeerStreamTrackState
  ) =>
    set((state) => {
      const usersPeersCopy = { ...state.usersPeers };
      usersPeersCopy[userId].isCameraEnabled = streamTrackState.camera;
      usersPeersCopy[userId].isMicEnabled = streamTrackState.mic;

      return {
        usersPeers: usersPeersCopy,
      };
    }),
}));
