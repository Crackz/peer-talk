"use client";
import { API_URL } from "@/constants";
import { AuthContext, UserInfo } from "@/modules/auth-provider";
import { useRouter } from "next/navigation";
import React, { useContext, useEffect, useState } from "react";

type Props = {};

const Page = (props: Props) => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const { authenticated } = useContext(AuthContext);

  const router = useRouter();

  useEffect(() => {
    if (authenticated) {
      router.push("/");
      return;
    }
  }, [authenticated]);

  const submitHandler = async (e: React.SyntheticEvent) => {
    e.preventDefault();
    try {
      const res = await fetch(`${API_URL}/v1/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });
      const data = await res.json();
      if (res.ok) {
        const user: UserInfo = {
          id: data.id,
          username: data.username,
          name: data.name,
        };
        localStorage.setItem("userInfo", JSON.stringify(data.user));
        localStorage.setItem("accessToken", data.accessToken);
        return router.push("/");
      }
    } catch (err) {
      console.log(err);
    }
  };

  return (
    <div className="flex items-center justify-center min-w-full min-h-screen">
      <form className="flex flex-col md:w-1/5">
        <div className="text-3xl font-bold text-center">
          <span className="text-white-500">welcome!</span>
        </div>
        <input
          placeholder="username"
          className="p-3 mt-8 rounded-md border-2 border-grey focus:outline-none focus:border-blue-500 focus:text-blue-500"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <input
          type="password"
          placeholder="password"
          className="p-3 mt-4 rounded-md border-2 border-grey focus:outline-none focus:border-blue-500 focus:text-blue-500"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <button
          className="p-3 mt-6 rounded-md bg-blue font-bold text-blue-500"
          type="submit"
          onClick={submitHandler}
        >
          login
        </button>
      </form>
    </div>
  );
};

export default Page;
