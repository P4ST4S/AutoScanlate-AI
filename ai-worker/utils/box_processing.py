"""Box processing utilities for text detection"""

from typing import List


def consolidate_boxes(boxes: List[List[int]], distance_threshold: int = 25) -> List[List[int]]:
    """
    Merge nearby text boxes that are close together.

    Args:
        boxes: List of bounding boxes [x1, y1, x2, y2]
        distance_threshold: Maximum distance between boxes to merge

    Returns:
        List of consolidated bounding boxes
    """
    if not boxes:
        return []

    rects = []
    for b in boxes:
        rects.append({
            'x1': b[0],
            'y1': b[1],
            'x2': b[2],
            'y2': b[3],
            'merged': False
        })

    merged = True
    while merged:
        merged = False
        new_rects = []
        while rects:
            current = rects.pop(0)
            i = 0
            while i < len(rects):
                other = rects[i]
                c_x1 = current['x1'] - distance_threshold
                c_y1 = current['y1'] - distance_threshold
                c_x2 = current['x2'] + distance_threshold
                c_y2 = current['y2'] + distance_threshold

                if (c_x1 < other['x2'] and c_x2 > other['x1'] and
                        c_y1 < other['y2'] and c_y2 > other['y1']):
                    new_box = {
                        'x1': min(current['x1'], other['x1']),
                        'y1': min(current['y1'], other['y1']),
                        'x2': max(current['x2'], other['x2']),
                        'y2': max(current['y2'], other['y2']),
                        'merged': False
                    }
                    current = new_box
                    rects.pop(i)
                    merged = True
                else:
                    i += 1
            new_rects.append(current)
        rects = new_rects

    return [[r['x1'], r['y1'], r['x2'], r['y2']] for r in rects]
