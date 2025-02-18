
import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { Layout } from "./components/layout/Layout";
import Dashboard from "./pages/Dashboard";
import StudyActivities from "./pages/StudyActivities";
import StudyActivity from "./pages/StudyActivity";
import Groups from "./pages/Groups";
import Group from "./pages/Group";
import StudySessions from "./pages/StudySessions";
import StudySession from "./pages/StudySession";
import Settings from "./pages/Settings";
import NotFound from "./pages/NotFound";

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <TooltipProvider>
      <Toaster />
      <Sonner />
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<Navigate to="/dashboard" replace />} />
            <Route path="dashboard" element={<Dashboard />} />
            <Route path="study_activities" element={<StudyActivities />} />
            <Route path="study_activity/:id" element={<StudyActivity />} />
            <Route path="groups" element={<Groups />} />
            <Route path="group/:id" element={<Group />} />
            <Route path="study_sessions" element={<StudySessions />} />
            <Route path="study_session/:id" element={<StudySession />} />
            <Route path="settings" element={<Settings />} />
            <Route path="*" element={<NotFound />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </TooltipProvider>
  </QueryClientProvider>
);

export default App;
