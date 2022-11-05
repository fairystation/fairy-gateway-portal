import './Profile.css';

import jwtDecode from 'jwt-decode';
import React, { useState, useEffect } from 'react';
import Blockies from 'react-blockies';
import food from './food.png'
import { Auth } from '../types';
import Web3 from 'web3';

interface Props {
	auth: Auth;
	onLoggedOut: () => void;
}

interface State {
	loading: boolean;
	user?: {
		id: number;
		username: string;
	};
	username: string;
	fake: boolean;
}

interface JwtDecoded {
	payload: {
		id: string;
		publicAddress: string;
	};
}

export const Profile = ({ auth, onLoggedOut }: Props): JSX.Element => {
	const [state, setState] = useState<State>({
		loading: false,
		user: undefined,
		username: '',
		fake: false
	});

	useEffect(() => {
		const { accessToken } = auth;
		const {
			payload: { id },
		} = jwtDecode<JwtDecoded>(accessToken);

		fetch(`${process.env.REACT_APP_BACKEND_URL}/users/${id}`, {
			headers: {
				Authorization: `Bearer ${accessToken}`,
			},
		})
			.then((response) => response.json())
			.then((user) => setState({ ...state, user }))
			.catch(window.alert);
	}, []);

	const handleChange = ({
		target: { value },
	}: React.ChangeEvent<HTMLInputElement>) => {
		setState({ ...state, username: value });
	};

	const handleConnectToService = async () => {
		const web3 = new Web3('http://localhost:8545');
		const contractAddress = '0xdeA7ED7764b919139A65a8A32FFFDC280F948e61';
		const contract = await import('./apps.json') as any;
		const myContract = new web3.eth.Contract(contract.abi, contractAddress);
		const result = await myContract.methods.mint('ass').call();
		console.log(result);
		setState({ ...state, fake: true });
		return result;
	};

	const handleSubmit = () => {
		const { accessToken } = auth;
		const { user, username } = state;

		setState({ ...state, loading: true });

		if (!user) {
			window.alert(
				'The user id has not been fetched yet. Please try again in 5 seconds.'
			);
			return;
		}

		fetch(`${process.env.REACT_APP_BACKEND_URL}/users/${user.id}`, {
			body: JSON.stringify({ username }),
			headers: {
				Authorization: `Bearer ${accessToken}`,
				'Content-Type': 'application/json',
			},
			method: 'PATCH',
		})
			.then((response) => response.json())
			.then((user) => setState({ ...state, loading: false, user }))
			.catch((err) => {
				window.alert(err);
				setState({ ...state, loading: false });
			});
	};

	const { accessToken } = auth;

	const {
		payload: { publicAddress },
	} = jwtDecode<JwtDecoded>(accessToken);

	const { loading, user } = state;

	const username = user && user.username;

	return (
		<div className="Profile">
			<p>
				Logged in as <Blockies seed={publicAddress} />
			</p>
			<div>
				My public address is: <pre>{publicAddress}</pre>
			</div>
			<p>
				{ 
				state.fake ? 
				<p> <h2>Your food is on the way!</h2>
				<br />
				<img src={food} alt="Tasty food" width="30%" className="Food-image" /></p>
				:
				<button onClick={handleConnectToService}>Connect to service</button>
				}
			</p> 
			<p>
				<button onClick={onLoggedOut}>Logout</button>
			</p>
		</div>
	);
};
