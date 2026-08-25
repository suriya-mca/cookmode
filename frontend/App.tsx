import { StatusBar } from "expo-status-bar";
import { StyleSheet, Text, View } from "react-native";

// TODO: API base URL — point at the Go backend (http://localhost:8080/api/v1 in dev).
export default function App() {
  return (
    <View style={styles.container}>
      <Text>CookMode</Text>
      <StatusBar style="auto" />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#fff",
    alignItems: "center",
    justifyContent: "center",
  },
});
