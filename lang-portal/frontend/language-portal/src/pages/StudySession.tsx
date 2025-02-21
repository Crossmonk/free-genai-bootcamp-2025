
import { useParams } from "react-router-dom";

const StudySession = () => {
  const { id } = useParams();
  return (
    <div className="animate-fade-in">
      <h1 className="text-3xl font-semibold mb-8">Study Session {id}</h1>
    </div>
  );
};

export default StudySession;
