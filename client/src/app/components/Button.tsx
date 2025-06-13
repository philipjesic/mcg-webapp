import { ButtonHTMLAttributes, FC } from "react";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "default" | "secondary";
  size?: "sm" | "md" | "lg";
}

export const Button: FC<ButtonProps> = ({
  children,
  className = "",
  variant = "default",
  size = "md",
  ...props
}) => {
  const baseStyles =
    "rounded px-4 py-2 font-semibold transition-colors duration-200 focus:outline-none focus:ring";

  const variantStyles =
    variant === "secondary"
      ? "bg-gray-200 text-black hover:bg-gray-300"
      : "bg-blue-600 text-white hover:bg-blue-700";

  const sizeStyles =
    size === "sm"
      ? "text-sm"
      : size === "lg"
      ? "text-lg py-3 px-5"
      : "text-base";

  const combinedClasses =
    `${baseStyles} ${variantStyles} ${sizeStyles} ${className}`.trim();

  return (
    <button className={combinedClasses} {...props}>
      {children}
    </button>
  );
};
