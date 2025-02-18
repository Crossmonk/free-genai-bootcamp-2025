
import { Bell } from "lucide-react";
import { Button } from "@/components/ui/button";

export const Header = () => {
  return (
    <header className="border-b">
      <div className="flex h-16 items-center px-4">
        <div className="ml-auto flex items-center space-x-4">
          <Button variant="ghost" size="icon">
            <Bell className="h-5 w-5" />
          </Button>
          <Button
            className="bg-primary hover:bg-primary-hover text-white"
            size="sm"
          >
            Launch Activity
          </Button>
        </div>
      </div>
    </header>
  );
};
