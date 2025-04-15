"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Spinner from "../components/Spinner";
import useRequest from "../hooks/use-request";

export default function Signup() {
  const router = useRouter();
  const [successMessage, setSuccessMessage] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const { doRequest, errors, loading } = useRequest({
    url: "/api/users/signup",
    method: "post",
    body: {
      email,
      password,
    },
    onSuccess: () => {
      setSuccessMessage("Sign up successful! Redirecting...");
      setTimeout(() => {
        router.push("/");
      }, 500);
    },
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    await doRequest();
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100 px-4">
      <form
        onSubmit={handleSubmit}
        className="bg-white p-8 rounded-2xl shadow-md w-full max-w-md space-y-6"
      >
        <h2 className="text-2xl font-semibold text-gray-800 text-center">
          Sign Up
        </h2>

        <div>
          <label
            htmlFor="email"
            className="block text-sm font-medium text-gray-700"
          >
            Email
          </label>
          <input
            name="email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-none"
          />
        </div>

        <div>
          <label
            htmlFor="password"
            className="block text-sm font-medium text-gray-700"
          >
            Password
          </label>
          <input
            name="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-none"
          />
        </div>

        <button
          type="submit"
          className={`w-full py-2 px-4 rounded-md text-white transition ${
            loading
              ? "bg-gray-400 cursor-not-allowed"
              : "bg-blue-600 hover:bg-blue-700"
          }`}
          disabled={loading || successMessage.length > 0}
        >
          {loading ? <Spinner message="Signing up"/> : "Create Account"}
        </button>
        {errors}

        {successMessage && <Spinner message={successMessage} />}
      </form>
    </div>
  );
}
