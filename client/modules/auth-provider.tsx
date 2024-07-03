"use client";
import { useState, createContext, useEffect } from "react";
import { useRouter } from "next/navigation";
import { AuthResponse } from "@/interfaces/auth-response";
import { USER_INFO_KEY } from "@/constants";

export type UserInfo = {
  id: string;
  name: string;
  username: string;
};

export const AuthContext = createContext<{
  authenticated: boolean;
  setAuthenticated: (auth: boolean) => void;
  user: UserInfo;
  setUser: (user: UserInfo) => void;
}>({
  authenticated: false,
  setAuthenticated: () => {},
  user: { username: "", name: "", id: "" },
  setUser: () => {},
});

const AuthContextProvider = ({ children }: { children: React.ReactNode }) => {
  const [authenticated, setAuthenticated] = useState(false);
  const [user, setUser] = useState<AuthResponse["user"]>({
    username: "",
    name: "",
    id: "",
  });

  const router = useRouter();

  useEffect(() => {
    const userInfoStr = localStorage.getItem(USER_INFO_KEY);

    if (!userInfoStr) {
      if (window.location.pathname != "/register") {
        router.push("/login");
        return;
      }
    } else {
      const userInfo: UserInfo = JSON.parse(userInfoStr);
      if (userInfo) {
        setUser(userInfo);
        setAuthenticated(true);
      }
    }
  }, [authenticated]);

  return (
    <AuthContext.Provider
      value={{
        authenticated: authenticated,
        setAuthenticated: setAuthenticated,
        user: user,
        setUser: setUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContextProvider;
