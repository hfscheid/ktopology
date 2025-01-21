import os
import pandas as pd
import torch
from torch_geometric.data import Data, DataLoader
from torch_geometric.nn import GCNConv
from torch.nn import Linear
from torch.nn.functional import cross_entropy

# Função para carregar os dados
def load_data(file_path):
    data = []
    with open(file_path, 'r') as f:
        for line in f:
            data.append(eval(line))
    return data

# Função para processar os dados em grafos
def process_data(data):
    nodes = []
    edges = []
    labels = []
    for entry in data:
        nodes.extend(entry['nodes'])
        edges.extend([(edge['source'], edge['target']) for edge in entry['edges']])
        labels.extend([node['metrics']['num_rejected'] for node in entry['nodes']])
    return nodes, edges, labels

# Definição do modelo GNN
class GNN(torch.nn.Module):
    def __init__(self, in_channels, hidden_channels, out_channels):
        super(GNN, self).__init__()
        self.conv1 = GCNConv(in_channels, hidden_channels)
        self.conv2 = GCNConv(hidden_channels, out_channels)
        self.lin = Linear(out_channels, 1)

    def forward(self, x, edge_index):
        x = self.conv1(x, edge_index).relu()
        x = self.conv2(x, edge_index).relu()
        x = self.lin(x)
        return x

# Função para treinar o modelo
def train(model, loader, optimizer, criterion):
    model.train()
    total_loss = 0
    for data in loader:
        optimizer.zero_grad()
        out = model(data.x, data.edge_index)
        loss = criterion(out, data.y)
        loss.backward()
        optimizer.step()
        total_loss += loss.item()
    return total_loss / len(loader)

# Função principal
def main():
    # Carregar e processar os dados
    file_path = '../sample-data/small-stable.jsonl'
    data = load_data(file_path)
    nodes, edges, labels = process_data(data)

    # Criar o grafo
    x = torch.tensor(nodes, dtype=torch.float)
    edge_index = torch.tensor(edges, dtype=torch.long).t().contiguous()
    y = torch.tensor(labels, dtype=torch.float)

    # Criar o objeto Data
    data = Data(x=x, edge_index=edge_index, y=y)

    # Criar o DataLoader
    loader = DataLoader([data], batch_size=1, shuffle=True)

    # Definir o modelo, otimizador e critério de perda
    model = GNN(in_channels=x.size(1), hidden_channels=16, out_channels=16)
    optimizer = torch.optim.Adam(model.parameters(), lr=0.01)
    criterion = cross_entropy

    # Treinar o modelo
    for epoch in range(100):
        loss = train(model, loader, optimizer, criterion)
        print(f'Epoch {epoch+1}, Loss: {loss:.4f}')

if __name__ == '__main__':
    main()