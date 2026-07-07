import { useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

function MeetingChatRedirect() {
  const { id: meetingId } = useParams();
  const navigate = useNavigate();
  const token = localStorage.getItem('token');
  const API_URL = import.meta.env.VITE_API_URL;

  useEffect(() => {
    const resolve = async () => {
      try {
        const res = await fetch(`${API_URL}/meetings/${meetingId}/chat`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        if (res.ok) {
          const { chat_id } = await res.json();
          navigate(`/group-chats/${chat_id}`, { replace: true });
        } else {
          navigate('/chats', { replace: true });
        }
      } catch {
        navigate('/chats', { replace: true });
      }
    };
    resolve();
  }, [meetingId, API_URL, token]);

  return <div className="rtr-loading">Загрузка...</div>;
}

export default MeetingChatRedirect;