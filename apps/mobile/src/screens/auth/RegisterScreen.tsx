import { useState } from 'react';
import { View } from 'react-native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { AppText, Button, Field, Screen } from '../../components';
import { apiErrorMessage } from '../../api/client';
import { useRegister } from '../../hooks/useAuth';
import { spacing } from '../../theme';
import type { AuthStackParamList } from '../../navigation/types';

type Props = NativeStackScreenProps<AuthStackParamList, 'Register'>;

export function RegisterScreen({ navigation }: Props) {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const registerMut = useRegister();

  const tooShort = password.length > 0 && password.length < 8;
  const canSubmit = !!name && !!email && password.length >= 8;

  const submit = () => {
    if (!canSubmit) return;
    registerMut.mutate({ name: name.trim(), email: email.trim().toLowerCase(), password });
  };

  return (
    <Screen scroll contentStyle={{ flexGrow: 1, justifyContent: 'center' }}>
      <View style={{ gap: spacing.xs, marginBottom: spacing.lg }}>
        <AppText variant="display">GET STARTED</AppText>
        <AppText variant="subtitle">Set up your training in a minute.</AppText>
      </View>

      <Field label="Name" value={name} onChangeText={setName} placeholder="Ali" autoCapitalize="words" />
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
        autoComplete="new-password"
        placeholder="At least 8 characters"
        error={tooShort ? 'Use at least 8 characters' : null}
        onSubmitEditing={submit}
        returnKeyType="go"
      />

      {registerMut.isError ? (
        <AppText variant="caption" color="#D64933">
          {apiErrorMessage(registerMut.error, 'Could not create account')}
        </AppText>
      ) : null}

      <Button
        label="Create account"
        fullWidth
        loading={registerMut.isPending}
        disabled={!canSubmit}
        onPress={submit}
        style={{ marginTop: spacing.sm }}
      />

      <View style={{ flexDirection: 'row', justifyContent: 'center', gap: spacing.xs, marginTop: spacing.md }}>
        <AppText variant="caption">Already have an account?</AppText>
        <AppText
          variant="caption"
          color="#E8A33D"
          onPress={() => navigation.navigate('Login')}
          suppressHighlighting
        >
          Log in
        </AppText>
      </View>
    </Screen>
  );
}
