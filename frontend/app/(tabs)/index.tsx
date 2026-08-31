import { View, Text } from "react-native";
import { theme } from "@/lib/theme";

export default function FeedScreen() {
  return (
    <View
      style={{
        flex: 1,
        backgroundColor: theme.colors.cream,
        alignItems: "center",
        justifyContent: "center",
        padding: theme.spacing.lg,
      }}
    >
      <Text
        style={{
          fontSize: theme.text.title,
          color: theme.colors.charcoal,
          fontWeight: "700",
        }}
      >
        Discover
      </Text>
      <Text
        style={{
          marginTop: 8,
          fontSize: theme.text.body,
          color: theme.colors.charcoalMuted,
        }}
      >
        Feed — GET /recipes + search chips coming next
      </Text>
    </View>
  );
}
