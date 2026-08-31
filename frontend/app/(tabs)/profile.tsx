import { View, Text } from "react-native";
import { theme } from "@/lib/theme";

export default function ProfileScreen() {
  return (
    <View
      style={{
        flex: 1,
        backgroundColor: theme.colors.cream,
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <Text style={{ color: theme.colors.charcoal, fontSize: theme.text.subtitle }}>
        Profile
      </Text>
      <Text style={{ color: theme.colors.charcoalMuted, marginTop: 8 }}>
        GET /users/:id/followers & following
      </Text>
    </View>
  );
}
