import axios, { Method, AxiosResponse } from "axios";
import { useState, ReactNode } from "react";

interface UserRequestProps {
  url: string;
  method: Method;
  body?: object;
  onSuccess?: (data: any) => void;
}

interface UserRequestReturn {
  doRequest: (props?: object) => Promise<any>;
  errors: ReactNode | null;
  loading: boolean;
}

const useRequest = ({
  url,
  method,
  body,
  onSuccess,
}: UserRequestProps): UserRequestReturn => {
  const [errors, setErrors] = useState<ReactNode | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const doRequest = async (props: object = {}) => {
    try {
      setLoading(true);
      setErrors(null);
      const response: AxiosResponse = await axios.request({
        method,
        url,
        data: { ...body, ...props },
      });

      if (onSuccess) {
        onSuccess(response.data);
      }

      return response.data;
    } catch (err: any) {
      if (err.response?.data?.errors) {
        console.log(err);
        setErrors(
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
            <strong className="font-bold">Oops!</strong>
            <ul className="mt-2 list-disc list-inside">
              {err.response.data.errors.map(
                (error: { message: string }, idx: number) => (
                  <li key={idx}>{error.message}</li>
                )
              )}
            </ul>
          </div>
        );
      }
    } finally {
      setLoading(false);
    }
  };

  return { doRequest, errors, loading };
};

export default useRequest;
