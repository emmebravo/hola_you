import { useState } from 'react';
import { useHMSActions } from '@100mslive/react-sdk';

const JoinRoom = () => {
  const hmsActions = useHMSActions();
  const [inputValues, setInputValues] = useState({
    name: '',
    token: '',
  });

  const handleInputChange = (event) => {
    const { name, value } = event.target;
    setInputValues((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    const { userName = '', roomCode = '' } = inputValues;

    // roomCode fetches auth token
    const authToken = await hmsActions.getAuthTokenByRoomCode({ roomCode });

    try {
      await hmsActions.join({ userName, authToken });
    } catch (error) {
      console.error(error);
    }
  };

  return <></>;
};

export default JoinRoom;
