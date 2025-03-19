// components/layout/Layout.tsx
import { NavLink, Outlet } from "react-router-dom";

export function Layout() {
  return (
    <div className="flex min-h-screen bg-white">
      <aside className="w-64 fixed h-screen bg-gray-100 p-4 border-r">
        <h2 className="text-xl font-bold text-blue-900 mb-6">Process Monitor</h2>
        <nav className="space-y-2">
          <NavLink 
            to="/" 
            className={({ isActive }) => 
              `block p-2 rounded ${isActive ? 'bg-blue-100 text-blue-800' : 'text-gray-700 hover:bg-gray-200'}`
            }
          >
            Process Logs
          </NavLink>
          <NavLink 
            to="/settings" 
            className={({ isActive }) => 
              `block p-2 rounded ${isActive ? 'bg-blue-100 text-blue-800' : 'text-gray-700 hover:bg-gray-200'}`
            }
          >
            Settings
          </NavLink>
        </nav>
      </aside>
      
      <main className="flex-1 ml-64 p-8">
        <Outlet />
      </main>
    </div>
  );
}