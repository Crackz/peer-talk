"use client";
import { ACCESS_TOKEN_KEY, USER_INFO_KEY } from "@/constants";
import { AuthResponse } from "@/interfaces/auth-response";
import { AuthContext } from "@/modules/auth-provider";
import axiosInstance from "@/utils/axios";
import axios from "axios";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useContext, useEffect, useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import { BounceLoader } from "react-spinners";

type Inputs = {
  name: string;
  username: string;
  email?: string;
  password: string;
};

const Page = () => {
  const router = useRouter();
  const { authenticated } = useContext(AuthContext);

  const [loading, setLoading] = useState(false);
  const [errorsMessages, setErrorsMessages] = useState<{
    errors?: { message: string; param: string }[];
  }>({});

  useEffect(() => {
    if (authenticated) {
      router.push("/");
      return;
    }
  }, [authenticated]);

  const {
    register,
    handleSubmit,
    formState: { errors: validationErrors },
  } = useForm<Inputs>();
  const onSubmit: SubmitHandler<Inputs> = async (data) => {
    console.log(data);
    setLoading(true);
    setErrorsMessages({});
    try {
      const response = await axiosInstance.post<AuthResponse>(
        "/v1/register",
        data
      );
      if (response.status === 200) {
        localStorage.setItem(ACCESS_TOKEN_KEY, response.data.accessToken);
        localStorage.setItem(USER_INFO_KEY, JSON.stringify(response.data.user));

        router.push("/lobby");
      }
    } catch (error) {
      if (axios.isAxiosError(error)) {
        setErrorsMessages(error.response?.data);
        return;
      }
      console.log("error", error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      className={
        "flex flex-col items-center justify-center h-screen p-5 w-full"
      }
    >
      <label className={"flex flex-col items-start"} htmlFor={"name"}>
        Name*
      </label>
      <input
        className={
          " text-black my-4 p-2 rounded border border-gray-300 focus:outline-none focus:ring-2 focus:ring-white"
        }
        type={"text"}
        id={"name"}
        placeholder={"name"}
        {...register("name", { required: true })}
      />{" "}
      {validationErrors.name && (
        <span className="mb-2 text-red-500">
          {validationErrors.name.message}
        </span>
      )}
      <label className={"flex flex-col items-start"} htmlFor={"username"}>
        Username*
      </label>
      <input
        className={
          " text-black my-4 p-2 rounded border border-gray-300 focus:outline-none focus:ring-2 focus:ring-white"
        }
        type={"text"}
        id={"username"}
        placeholder={"Username"}
        {...register("username", { required: true })}
      />
      {validationErrors.username && (
        <span className="mb-2 text-red-500">
          {validationErrors.username.message}
        </span>
      )}
      <label htmlFor={"email"}>Email</label>
      <input
        className={
          "text-black my-4 p-2 rounded border border-gray-300 focus:outline-none focus:ring-2 focus:ring-white"
        }
        type={"text"}
        id={"email"}
        placeholder={"Email"}
        {...register("email", {
          pattern: {
            value: /\S+@\S+\.\S+/,
            message: "Entered value does not match email format",
          },
        })}
      />
      {validationErrors.email && (
        <span className="mb-2 text-red-500">
          {validationErrors.email.message}
        </span>
      )}
      <label htmlFor={"password"}>Password*</label>
      <input
        className={
          "text-black my-4 p-2 rounded border border-gray-300 focus:outline-none focus:ring-2 focus:ring-white"
        }
        type={"password"}
        id={"password"}
        placeholder={"Password"}
        {...register("password", {
          required: true,
          minLength: {
            value: 4,
            message: "Password must be at least 4 characters",
          },
        })}
      />
      {validationErrors.password && (
        <span className="mb-2 text-red-500">
          {validationErrors.password.message}
        </span>
      )}
      <button
        className={
          "action__button__container action__button2 cursor-pointer my-4"
        }
        onClick={handleSubmit(onSubmit)}
        disabled={loading}
      >
        {loading ? <BounceLoader size={20} /> : "Sign Up"}
      </button>
      {errorsMessages.errors &&
        errorsMessages.errors.map((errMsg, idx) => (
          <span key={idx} className="mb-2 text-red-500">
            {errMsg.message}
          </span>
        ))}
      <Link className={"text-sm my-4"} href={"/login"}>
        Already have an account? Login
      </Link>
    </div>
  );
};
export default Page;
