import { Navigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

const ProtectedRoute = ({ children }) => {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return <div>Загрузка...</div>;
  }

  // if (!isAuthenticated) {
  //   return <Navigate to="/welcome" replace />;
  // }

  return children;
};

export default ProtectedRoute;