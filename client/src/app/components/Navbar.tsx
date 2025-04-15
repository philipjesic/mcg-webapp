import Link from "next/link";

export default function Navbar() {
  return (
    <nav className="bg-white shadow-md sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16 items-center">
          {/* Left - Logo */}
          <div className="flex-shrink-0">
            <Link href="/">
              <span className="text-xl font-bold text-blue-600">MCG</span>
            </Link>
          </div>

          {/* Center - Navigation links (optional) */}
          <div className="hidden md:flex space-x-8">
            <Link href="/about" className="text-gray-700 hover:text-blue-600">
              About
            </Link>
            <Link
              href="/features"
              className="text-gray-700 hover:text-blue-600"
            >
              Features
            </Link>
          </div>

          {/* Right - Sign In */}
          <div className="flex items-center space-x-4">
            <Link
              href="/signin"
              className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded hover:bg-blue-700 transition"
            >
              Sign In
            </Link>
          </div>
        </div>
      </div>
    </nav>
  );
}
