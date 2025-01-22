import os
import json
import numpy as np
import pandas as pd
import tensorflow as tf
from spektral.data import Dataset, Graph
from spektral.data.loaders import SingleLoader
from spektral.layers import GCNConv
from tensorflow.keras.models import Model
from tensorflow.keras.layers import Dense
from tensorflow.keras.optimizers import Adam

# Função para carregar os dados
def load_data(file_path):
    data = []
    with open(file_path, 'r') as f:
        for line in f:
            data.append(json.loads(line))
    return data

# Função para processar os dados em grafos
def process_data(data):
    nodes = []
    edges = []
    labels = []
    for entry in data:
        nodes.extend(entry['nodes'])
        edges.extend([(edge['source'], edge['target']) for edge in entry['edges']])
        labels.extend([node['metrics'].get('num_rejected', 0) for node in entry['nodes']])
    return nodes, edges, labels

# Definição do Dataset
class MyDataset(Dataset):
    def __init__(self, file_path, **kwargs):
        self.file_path = file_path
        super().__init__(**kwargs)

    def read(self):
        data = load_data(self.file_path)
        nodes, edges, labels = process_data(data)
        x = np.array(nodes, dtype=np.float32)
        a = np.array(edges, dtype=np.int64).T
        y = np.array(labels, dtype=np.float32)
        return [Graph(x=x, a=a, y=y)]

# Definição do modelo GNN
class GNN(Model):
    def __init__(self):
        super().__init__()
        self.conv1 = GCNConv(16, activation='relu')
        self.conv2 = GCNConv(16, activation='relu')
        self.dense = Dense(1)

    def call(self, inputs):
        x, a = inputs
        x = self.conv1([x, a])
        x = self.conv2([x, a])
        return self.dense(x)

# Função principal
def main():
    # Carregar e processar os dados
    file_path = 'sample-data/small-stable.jsonl'
    dataset = MyDataset(file_path)
    loader = SingleLoader(dataset)

    print('File path ', file_path)
    print('Dataset ', dataset)
    print('Loader', loader)

    # Definir o modelo, otimizador e critério de perda
    model = GNN()
    optimizer = Adam(learning_rate=0.01)
    loss_fn = tf.keras.losses.MeanSquaredError()

    # Treinar o modelo
    for epoch in range(100):
        for batch in loader:
            with tf.GradientTape() as tape:
                inputs = batch[0], batch[1]
                targets = batch[2]
                predictions = model(inputs, training=True)
                loss = loss_fn(targets, predictions)
            gradients = tape.gradient(loss, model.trainable_variables)
            optimizer.apply_gradients(zip(gradients, model.trainable_variables))
        print(f'Epoch {epoch+1}, Loss: {loss.numpy():.4f}')

if __name__ == '__main__':
    main()