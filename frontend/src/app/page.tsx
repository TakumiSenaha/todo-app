import Link from "next/link";

export default function Home() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full space-y-8 text-center">
        <div>
          <h1 className="text-4xl font-bold text-gray-900">Todo App</h1>
          <p className="mt-4 text-lg text-gray-600">
            A full-stack application with authentication
          </p>
        </div>

        <div className="space-y-4">
          <Link
            href="/login"
            className="w-full flex justify-center py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            Sign In
          </Link>

          <Link
            href="/register"
            className="w-full flex justify-center py-3 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            Create Account
          </Link>
        </div>

        <div className="mt-8 text-sm text-gray-500">
          <p>Built with:</p>
          <ul className="mt-2 space-y-1">
            <li>• Go (Backend API)</li>
            <li>• Next.js (Frontend)</li>
            <li>• JWT Authentication</li>
            <li>• PostgreSQL Database</li>
          </ul>
        </div>
      </div>
    </div>
  );
}
