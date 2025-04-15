"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import useRequest from "../hooks/use-request";
import Spinner from "../components/Spinner";
import Link from "next/link";

export default function SignInPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const { doRequest, errors, loading } = useRequest({
    url: "/api/users/login",
    method: "post",
    body: {
      email,
      password,
    },
    onSuccess: () => {
      setSuccessMessage("Sign in successful! Redirecting...");
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
      <div className="bg-white shadow-md rounded-lg p-8 w-full max-w-md">
        <form onSubmit={handleSubmit} className="space-y-4">
        <h2 className="text-2xl font-bold mb-6 text-center text-gray-800">
          Sign in to your account
        </h2>
          <div>
            <label
              htmlFor="email"
              className="block text-sm font-medium text-gray-700"
            >
              Email
            </label>
            <input
              id="email"
              type="email"
              required
              className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-none"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
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
              id="password"
              type="password"
              required
              className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-none"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
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
            {loading ? <Spinner message="Signing In..." /> : "Sign In"}
          </button>

          <p className="text-sm text-center text-gray-500 mt-4">
            Don't have an account?{" "}
            <Link href="/signup" className="text-blue-600 hover:underline">
              Sign up
            </Link>
          </p>
          {errors}
          {successMessage && <Spinner message={successMessage} />}
        </form>
      </div>
    </div>
  );
}
