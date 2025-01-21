import json
import tensorflow as tf
import tensorflow_gnn as tfgnn

sample_graph = [
  {
    "nodes": [
      {"id": 0, "feature": [1.0, 2.0]},
      {"id": 1, "feature": [2.0, 3.0]},
      {"id": 2, "feature": [3.0, 4.0]}
    ],
    "edges": [
      {"source": 0, "target": 1},
      {"source": 1, "target": 2}
    ]
  },
  {
    "nodes": [
      {"id": 0, "feature": [4.0, 5.0]},
      {"id": 1, "feature": [5.0, 6.0]},
      {"id": 2, "feature": [6.0, 7.0]}
    ],
    "edges": [
      {"source": 0, "target": 2},
      {"source": 2, "target": 1}
    ]
  }
]

# Hilfsfunktion: Graph aus JSON-Daten erstellen
def create_graph_from_json(graph_data):
    nodes = graph_data["nodes"]
    edges = graph_data["edges"]

    # Anzahl der Knoten und Kanten
    num_nodes = len(nodes)
    num_edges = len(edges)

    # Node-Features und IDs extrahieren
    node_features = [node["feature"] for node in nodes]
    node_ids = [node["id"] for node in nodes]

    # Edges (Source- und Target-Knoten)
    source_nodes = [edge["source"] for edge in edges]
    target_nodes = [edge["target"] for edge in edges]

    # GraphTensor erstellen
    graph_tensor = tfgnn.GraphTensor.from_pieces(
        node_sets={
            "nodes": tfgnn.NodeSet.from_fields(
                sizes=tf.constant([num_nodes]),
                features={"feature": tf.constant(node_features)}
            )
        },
        edge_sets={
            "edges": tfgnn.EdgeSet.from_fields(
                sizes=tf.constant([num_edges]),
                adjacency=tfgnn.Adjacency.from_indices(
                    source=("nodes", tf.constant(source_nodes)),
                    target=("nodes", tf.constant(target_nodes))
                )
            )
        }
    )
    return graph_tensor

# JSON-Datei mit mehreren Graphen laden
with open("graphs.json", "r") as file:
    graph_list = json.load(file)

# GraphTensor für jeden Graphen erstellen
graph_tensors = [create_graph_from_json(graph) for graph in graph_list]

# Beispiel: Alle Graphen anzeigen
for i, graph in enumerate(graph_tensors):
    print(f"Graph {i+1}:")
    print(graph)


# Dataset aus GraphTensors erstellen
dataset = tf.data.Dataset.from_tensor_slices(graph_tensors)

# Beispiel: Dataset für Training vorbereiten
def preprocess(graph):
    # Extrahiere Features und Labels (falls vorhanden)
    node_features = graph.node_sets["nodes"]["feature"]
    labels = tf.constant([1])  # Beispiel-Label
    return node_features, labels

dataset = dataset.map(preprocess).batch(2)  # Batches erstellen

# Beispiel-Training (mit einfachem Dummy-Modell)
model = tf.keras.Sequential([
    tf.keras.layers.Input(shape=(2,)),  # Beispiel-Featuregröße: 2
    tf.keras.layers.Dense(4, activation="relu"),
    tf.keras.layers.Dense(1, activation="sigmoid")
])

# Kompilieren und trainieren
model.compile(optimizer="adam", loss="binary_crossentropy")
model.fit(dataset, epochs=10)
