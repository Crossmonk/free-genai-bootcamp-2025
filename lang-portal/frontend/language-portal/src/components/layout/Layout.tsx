
import { Outlet } from "react-router-dom";
import { Sidebar } from "./Sidebar";
import { Header } from "./Header";

export const Layout = () => {
  return (
    <div className="min-h-screen bg-background">
      <div className="flex h-screen overflow-hidden">
        <Sidebar />
        <div className="flex-1 overflow-auto">
          <Header />
          <main className="container mx-auto py-6 px-4 animated fade-in">
            <Outlet />
          </main>
        </div>
      </div>
    </div>
  );
};
