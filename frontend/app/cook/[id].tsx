import { View, Text } from "react-native";
import { useLocalSearchParams } from "expo-router";
import { theme } from "@/lib/theme";

export default function CookMode() {
  const { id } = useLocalSearchParams<{ id: string }>();
  return (
    <View
      style={{
        flex: 1,
        backgroundColor: theme.colors.charcoal,
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <Text style={{ color: theme.colors.cream, fontSize: theme.text.title }}>
        Cook Mode — {id}
      </Text>
      <Text style={{ color: theme.colors.cream, marginTop: 8, opacity: 0.7 }}>
        Video 60% + bottom sheet steps, anchor_seconds seek
      </Text>
    </View>
  );
}
