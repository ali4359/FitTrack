import { useState } from 'react';
import { View } from 'react-native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { AppText, Button, Field, Screen } from '../../components';
import { apiErrorMessage } from '../../api/client';
import { useLogin } from '../../hooks/useAuth';
import { spacing } from '../../theme';
import type { AuthStackParamList } from '../../navigation/types';

type Props = NativeStackScreenProps<AuthStackParamList, 'Login'>;

export function LoginScreen({ navigation }: Props) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const login = useLogin();

  const submit = () => {
    login.mutate({ email: email.trim().toLowerCase(), password });
  };

  return (
    <Screen scroll contentStyle={{ flexGrow: 1, justifyContent: 'center' }}>
      <View style={{ gap: spacing.xs, marginBottom: spacing.lg }}>
        <AppText variant="display">IRON &amp; SPICE</AppText>
        <AppText variant="subtitle">Train hard. Eat right. Locally.</AppText>
      </View>

      <Field
        label="Email"
        value={email}
        onChangeText={setEmail}
        autoCapitalize="none"
        keyboardType="email-address"
        autoComplete="email"
        placeholder="you@example.com"
      />
      <Field
        label="Password"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
        autoComplete="current-password"
        placeholder="••••••••"
        onSubmitEditing={submit}
        returnKeyType="go"
      />

      {login.isError ? (
        <AppText variant="caption" color="#D64933">
          {apiErrorMessage(login.error, 'Could not sign in')}
        </AppText>
      ) : null}

      <Button
        label="Log in"
        fullWidth
        loading={login.isPending}
        disabled={!email || !password}
        onPress={submit}
        style={{ marginTop: spacing.sm }}
      />

      <View style={{ flexDirection: 'row', justifyContent: 'center', gap: spacing.xs, marginTop: spacing.md }}>
        <AppText variant="caption">New here?</AppText>
        <AppText
          variant="caption"
          color="#E8A33D"
          onPress={() => navigation.navigate('Register')}
          suppressHighlighting
        >
          Create an account
        </AppText>
      </View>
    </Screen>
  );
}
