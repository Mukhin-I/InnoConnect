import { useState, useEffect } from 'react';

export const useAuth = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('token');
    
    if (!token) {
      setIsAuthenticated(false);
      setIsLoading(false);
      return;
    }
    
    try {
      //const payload = JSON.parse(atob(token.split('.')[1]));
      const base64Url = token.split('.')[1];

      const base64 = base64Url
        .replace(/-/g, '+')
        .replace(/_/g, '/')
        .padEnd(base64Url.length + (4 - base64Url.length % 4) % 4, '=');

      const payload = JSON.parse(atob(base64));

      const exp = payload.exp * 1000;
      const currentTime = Date.now();
      
      if (exp > currentTime) {
        setIsAuthenticated(true);
      } else {
        localStorage.removeItem('token');
        localStorage.removeItem('tokenType');
        localStorage.removeItem('expiresIn');
        setIsAuthenticated(false);
      }
    } catch (error) {
      console.error('Invalid token:', error);
      localStorage.removeItem('token');
      localStorage.removeItem('tokenType');
      localStorage.removeItem('expiresIn');
      setIsAuthenticated(false);
    }
    
    setIsLoading(false);
  }, []);

  return { isAuthenticated, isLoading };
};