import { useAuth } from '~/hooks';

const Restricted: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user } = useAuth();
  if (!user || user.role == 'student') {
    return null;
  }
  if (user.role == 'admin' || user.role == 'teacher') {
    return <>{children}</>;
  }
};

export default Restricted;
