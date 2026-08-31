import { View, Text } from "react-native";
import { theme } from "@/lib/theme";

export default function CollectionsScreen() {
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
        Collections
      </Text>
      <Text style={{ color: theme.colors.charcoalMuted, marginTop: 8 }}>
        GET /collections + /collections/:id/recipes
      </Text>
    </View>
  );
}
